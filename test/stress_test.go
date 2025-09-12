package test

import (
	"PBL/client/utils"
	"PBL/shared"
	"fmt"
	"io/fs"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type TestRoom struct {
	Players []string
}

type ConcurrencyMetrics struct{
	TotalConnections int64
	SuccessfulLogins int64
	FailedLogins int64
	PlayersInQueue int64
	MatchesCreated int64
	PacksOpened int64
	TotalMessages int64 
	CardsDistributed map[string]int64
	UniqueCards map[string]bool
	Latence []time.Duration
	mutex sync.RWMutex
}

var concurrencyMetrics = ConcurrencyMetrics{
	CardsDistributed: make(map[string]int64),
	UniqueCards:      make(map[string]bool),
}

var (
	MatchedRooms = make(map[string]TestRoom)
	roomMutex    sync.Mutex
	testWG       sync.WaitGroup
)

const numClient = 200

//Teste para a concorrência no momento do login -> há a lista de usuários online
func TestConcurrentLogin(t *testing.T){
	defer CleanTestUser() //limpar os arquivos teste
	testWG.Add(numClient) //numro de goroutines que vão rodar
	startTime := time.Now()
	for i := 1; i <= numClient; i++{
		go testLogin(i)
		time.Sleep(1 * time.Microsecond)
	}
	testWG.Wait()
	fmt.Println("Tempo total: ", time.Since(startTime))
	fmt.Println("Conexões: ", atomic.LoadInt64(&concurrencyMetrics.TotalConnections))
	fmt.Println("Logins falhados: ", atomic.LoadInt64(&concurrencyMetrics.FailedLogins))
	
}

func TestConcurrentPack(t *testing.T){
	defer CleanTestUser() //limpar os arquivos teste
	testWG.Add(numClient)
	startTime := time.Now()
	for i := 1; i <= numClient; i++{
		go testOpenPack(i)
		time.Sleep(1 * time.Microsecond)
	}
	testWG.Wait()
	fmt.Println("Tempo total: ", time.Since(startTime))
	fmt.Printf("Pacotes abertos: %d\n", atomic.LoadInt64(&concurrencyMetrics.PacksOpened))
}

//Para tentar a conexão algumas vezes e não só uma
func connectWithRetry(id int) (net.Conn, error) {
	for i := 0; i < 3; i++ {
		conn, err := net.Dial("tcp", "localhost:8080")
		if err == nil {
			return conn, nil
		}
		time.Sleep(10 * time.Millisecond)
	}
	fmt.Printf("Cliente %d falhou ao conectar após 3 tentativas\n", id)
	return nil, fmt.Errorf("falha na conexão")
}

//LOGIN
func testLogin(id int) {
	defer testWG.Done()
	conn, err := connectWithRetry(id)
	if err != nil {
		atomic.AddInt64(&concurrencyMetrics.FailedLogins, 1)
		atomic.AddInt64(&concurrencyMetrics.TotalConnections, 1)
		return
	}
	defer conn.Close()

	atomic.AddInt64(&concurrencyMetrics.TotalConnections, 1)

	respChan := make(chan shared.Response, 5)
	stopChan := make(chan bool)
	go utils.ListenServer(conn, respChan, stopChan)

	//Register
	user := shared.User{
		UserName: fmt.Sprintf("testUser%d", id),
		Password: "123",
	}
	
	utils.SendRequest(conn, "REGISTER", user)
	time.Sleep(10 * time.Microsecond)
	
	utils.SendRequest(conn, "LOGIN", user)
	
	timeout := time.After(3 * time.Second)
	loginSuccess := false
	
	for {
		select {
		case resp := <-respChan:
			atomic.AddInt64(&concurrencyMetrics.TotalMessages, 1)
			if resp.Status == "successLogin" {
				atomic.AddInt64(&concurrencyMetrics.SuccessfulLogins, 1)
				loginSuccess = true
				return
			}
		case <-timeout:
			if !loginSuccess {
				atomic.AddInt64(&concurrencyMetrics.FailedLogins, 1)
			}
			return
		}
	}
}

//PACK
func testOpenPack(id int) {
	defer testWG.Done()
	conn, err := connectWithRetry(id)
	if err != nil {
		fmt.Printf("Cliente %d não conseguiu conectar\n", id)
		return
	}
	defer conn.Close()

	respChan := make(chan shared.Response, 5)
	stopChan := make(chan bool)
	go utils.ListenServer(conn, respChan, stopChan)

	user := shared.User{
		UserName: fmt.Sprintf("testUser%d", id),
		Password: "123",
	}

	utils.SendRequest(conn, "REGISTER", user)
	time.Sleep(10 * time.Microsecond)
	utils.SendRequest(conn, "LOGIN", user)
	time.Sleep(10 * time.Microsecond)
	utils.SendRequest(conn, "PACK", user)

	timeout := time.After(5 * time.Second)
	for {
		select {
		case resp := <-respChan:
			atomic.AddInt64(&concurrencyMetrics.TotalMessages, 1)

			if resp.Status == "successPack" {
				atomic.AddInt64(&concurrencyMetrics.PacksOpened, 1)

				//Extrair cartas da resposta
				if data, ok := resp.Data.(map[string]interface{}); ok {
					if cards, ok := data["cards"].([]interface{}); ok {
						concurrencyMetrics.mutex.Lock()
						for _, card := range cards {
							cardStr := card.(string)

							if concurrencyMetrics.UniqueCards[cardStr] {
								fmt.Printf("ERRO: carta duplicada detectada! %s\n", cardStr)
							} else {
								concurrencyMetrics.UniqueCards[cardStr] = true
								concurrencyMetrics.CardsDistributed[cardStr]++
							}
						}
						concurrencyMetrics.mutex.Unlock()
					}
				}

				fmt.Printf("Cliente %d abriu pacote com sucesso!\n", id)
				return
			}
		case <-timeout:
			fmt.Printf("Cliente %d timeout ao abrir pacote\n", id)
			return
		}
	}
}

//Apagar arquivos json dos usuários teste
func CleanTestUser() error {
    dir := "../server/data" 
    return filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() && strings.HasPrefix(info.Name(), "test") && strings.HasSuffix(info.Name(), ".json") {
            os.Remove(path)
        }
        return nil
    })
}
