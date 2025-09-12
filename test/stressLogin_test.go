package main

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

type LoginMetrics struct {
	TotalConnections  int64
	SuccessfulLogins  int64
	FailedLogins      int64
	TotalMessages     int64
}

var loginMetrics = LoginMetrics{}

var (
	MatchedRooms = make(map[string]TestRoom)
	testWG       sync.WaitGroup
)

const loginNumClient = 200

//Teste para a concorrência no momento do login -> há a lista de usuários online
func TestConcurrentLogin(t *testing.T){
	defer CleanTestUser() //limpar os arquivos teste
	testWG.Add(loginNumClient) //numro de goroutines que vão rodar
	startTime := time.Now()
	for i := 1; i <= loginNumClient; i++{
		go testLogin(i)
		time.Sleep(1 * time.Microsecond)
	}
	testWG.Wait()
	fmt.Println("Tempo total: ", time.Since(startTime))
	fmt.Println("Conexões: ", atomic.LoadInt64(&loginMetrics.TotalConnections))
	fmt.Println("Logins falhados: ", atomic.LoadInt64(&loginMetrics.FailedLogins))
	
}

//Para tentar a conexão algumas vezes e não só uma
func ConnectWithRetry(id int) (net.Conn, error) {
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
	conn, err := ConnectWithRetry(id)
	if err != nil {
		atomic.AddInt64(&loginMetrics.FailedLogins, 1)
		atomic.AddInt64(&loginMetrics.TotalConnections, 1)
		return
	}
	defer conn.Close()

	atomic.AddInt64(&loginMetrics.TotalConnections, 1)

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
			atomic.AddInt64(&loginMetrics.TotalMessages, 1)
			if resp.Status == "successLogin" {
				atomic.AddInt64(&loginMetrics.SuccessfulLogins, 1)
				loginSuccess = true
				return
			}
		case <-timeout:
			if !loginSuccess {
				atomic.AddInt64(&loginMetrics.FailedLogins, 1)
			}
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