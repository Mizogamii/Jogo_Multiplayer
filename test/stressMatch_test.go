package main

import (
	"PBL/client/utils"
	"PBL/shared"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"
)

type Metrics struct {
	Latencies []time.Duration
	QueueTimes []time.Duration
	Messages int
	mutex sync.Mutex
}

var metrics = Metrics{}

var (
	CreateRoomsTest = make(map[string]TestRoom) 
	QueueTest []string
	mu sync.Mutex
	waitG sync.WaitGroup
)

const matchNumClient = 200

func TestConcurrentMatchmaking(t *testing.T){
	defer CleanTestUser()
	waitG.Add(matchNumClient)

	for i := 1; i <= matchNumClient; i++ {
		go func(id int) {
			defer waitG.Done()
			testMatchmaking(id)
		}(i)
	}
	waitG.Wait()
	reportMatchmaking()
}

func testMatchmaking(id int){
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Printf("Cliente %d não conseguiu conectar: %v\n", id, err)
		return
	}
	defer conn.Close()

	respChan := make(chan shared.Response, 10)
	stopChan := make(chan bool)
	go utils.ListenServer(conn, respChan, stopChan)

	user := shared.User{
		UserName: fmt.Sprintf("teste%d", id),
		Password: "123",
	}
	utils.SendRequest(conn, "REGISTER", user)
	time.Sleep(10 * time.Millisecond)

	utils.SendRequest(conn, "LOGIN", user)
	time.Sleep(10 * time.Millisecond)

	simulateEnterQueue(conn, id)
	
	baseTimeout := 30 * time.Second
	if matchNumClient > 1000 {
		baseTimeout = time.Duration(matchNumClient/50) * time.Second
	}

	if baseTimeout > 2*time.Minute {
		baseTimeout = 2 * time.Minute
	}
	
	timeout := time.After(baseTimeout)

	for {
		select {
		case resp, ok := <-respChan:
			if !ok {
				return
			}

			metrics.mutex.Lock()
			metrics.Messages++
			metrics.mutex.Unlock()
			start := time.Now() 

			switch resp.Status {
			case "match":
				handleMatch(resp, id)
				elapsed := time.Since(start)
				metrics.mutex.Lock()
				metrics.Latencies = append(metrics.Latencies, elapsed)
				metrics.mutex.Unlock()
				return

			case "successPlay":
				fmt.Printf("Cliente %d na fila...\n", id)
				
			default:
				fmt.Printf("Cliente %d recebeu: %+v\n", id, resp)
			}

		case <-timeout:
			fmt.Printf("Cliente %d timeout, encerrando\n", id)
			return
		}
	}
}

func simulateEnterQueue(conn net.Conn, id int) {
	user := shared.User{
		UserName: fmt.Sprintf("teste%d", id),
		Password: "123",
	}
	err := utils.SendRequest(conn, "PLAY", user)
	if err != nil {
		fmt.Printf("Cliente %d erro ao entrar na fila: %v\n", id, err)
		return
	}
}

func handleMatch(resp shared.Response, id int) {
	player1 := fmt.Sprintf("teste%d", id)
	player2 := fmt.Sprintf("%v", resp.Data)

	if player1 > player2 {
		player1, player2 = player2, player1
	}
	key := fmt.Sprintf("%s|%s", player1, player2)

	mu.Lock()
	defer mu.Unlock()
	if _, exists := CreateRoomsTest[key]; !exists {
		CreateRoomsTest[key] = TestRoom{Players: []string{player1, player2}}
	}
	fmt.Printf("Cliente %d entrou em partida: %s\n", id, key)
}

func reportMatchmaking() {
    mu.Lock()
    defer mu.Unlock()
	fmt.Println("\n----------------------------------")
    fmt.Println("Salas criadas:")
    i := 1
    for _, room := range CreateRoomsTest {
        fmt.Printf("Sala %d: %v\n", i, room.Players)
        i++
    }

    //Relatório de latência
    metrics.mutex.Lock()
    defer metrics.mutex.Unlock()
    var total time.Duration
    for _, l := range metrics.Latencies {
        total += l
    }

    avgLatency := time.Duration(0)
    if len(metrics.Latencies) > 0 {
        avgLatency = total / time.Duration(len(metrics.Latencies))
    }
	fmt.Println("\n----------------------------------")
    fmt.Println("             Latências            ")
	fmt.Println("----------------------------------")
    fmt.Printf("Mensagens recebidas: %d\n", metrics.Messages)
    fmt.Printf("Latência média: %v\n", avgLatency)
    fmt.Println("=====================")
}
