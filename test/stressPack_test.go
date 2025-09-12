package main

import (
	"PBL/client/utils"
	"PBL/shared"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var (
	cartasDistribuidas = make(map[string]string) 
	mutex              sync.Mutex
	wg                 sync.WaitGroup
	duplicatas         int64
	iniciais           = map[string]bool{
		"AGUA": true, "TERRA": true, "FOGO": true, "AR": true, "MATO": true,
	}
)
	
var packNumClient = 200

func TestConcurrentPack(t *testing.T) {
	defer CleanTestUser() //limpar os arquivos teste
	wg.Add(packNumClient)

	for i := 1; i <= packNumClient; i++ {
		go func(id int) {
			defer wg.Done()
			testOpenPack(id)
		}(i)
	}

	wg.Wait()

	totalDuplicatas := atomic.LoadInt64(&duplicatas)
	if totalDuplicatas > 0 {
		t.Errorf("Encontradas %d cartas duplicadas!", totalDuplicatas)
	} else {
		fmt.Println("Nenhuma carta duplicada encontrada!")
	}
}

//Teste abrindo os pacotes
func testOpenPack(id int) {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Printf("Cliente %d não conseguiu conectar\n", id)
		return
	}
	defer conn.Close()

	user := shared.User{
		UserName: fmt.Sprintf("testUser%d", id),
		Password: "123",
	}

	utils.SendRequest(conn, "REGISTER", user)
	time.Sleep(10 * time.Millisecond)
	utils.SendRequest(conn, "LOGIN", user)
	time.Sleep(10 * time.Millisecond)

	utils.SendRequest(conn, "PACK", user)

	respChan := make(chan shared.Response, 5)
	stopChan := make(chan bool)
	go utils.ListenServer(conn, respChan, stopChan)

	timeout := time.After(5 * time.Second)
	for {
		select {
		case resp := <-respChan:
			if resp.Status == "successPack" {
				if data, ok := resp.Data.(map[string]interface{}); ok {
					if cards, ok := data["cards"].([]interface{}); ok {
						mutex.Lock()
						for _, c := range cards {
							carta := c.(string)

							//Ignora cartas iniciais -> (todos os usuários iniciam com 5 cartas padrões então todos tem cartas repetidas)
							if iniciais[carta] {
								continue
							}

							if dono, existe := cartasDistribuidas[carta]; existe {
								fmt.Printf("DUPLICATA! Carta %s já foi para %s\n", carta, dono)
								atomic.AddInt64(&duplicatas, 1)
							} else {
								cartasDistribuidas[carta] = user.UserName
							}
						}
						mutex.Unlock()
					}
				}
				return
			} else if resp.Status == "error" && resp.Data == "acabou as cartas" {
				fmt.Printf("Cliente %d: cartas esgotadas\n", id)
				return
			}
		case <-timeout:
			fmt.Printf("Cliente %d: timeout\n", id)
			return
		}
	}
}

