package main

import (
	"PBL/client/utils"
	"PBL/shared"
	"fmt"
	"net"
	"sync"
	"time"
)

type TestRoom struct {
	Players []string
}

var (
	CreateRoomsTest []TestRoom
	QueueTest       []string
	mu              sync.Mutex
	wg              sync.WaitGroup
)

func main() {
	numClients := 5
	wg.Add(numClients) 

	for i := 1; i <= numClients; i++ {
		go simulateClient(i)
	}

	wg.Wait()

	// Relatório final
	fmt.Println("----- RELATÓRIO FINAL -----")
	fmt.Println("Salas criadas:")
	for i, sala := range CreateRoomsTest {
		fmt.Printf("Sala %d: %v\n", i+1, sala.Players)
	}
	fmt.Println("Clientes ainda na fila:", QueueTest)
	fmt.Println("---------------------------")
}

func simulateClient(id int) {
	defer wg.Done()

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Printf("Cliente %d: erro ao conectar: %v\n", id, err)
		return
	}
	defer conn.Close()

	fmt.Printf("Cliente %d conectado\n", id)

	respChan := make(chan shared.Response)
	stopChan := make(chan bool)

	go utils.ListenServer(conn, respChan, stopChan)

	simulateRegister(conn, id)
	time.Sleep(50 * time.Millisecond)
	simulateLogin(conn, id)
	time.Sleep(50 * time.Millisecond)
	simulateEnterQueue(conn, id)

	for resp := range respChan {
		mu.Lock()
		switch resp.Status {
		case "match":
			CreateRoomsTest = append(CreateRoomsTest, TestRoom{
				Players: []string{idToName(id), resp.Data.(string)},
			})
		case "successPlay":
			QueueTest = append(QueueTest, idToName(id))
		}
		mu.Unlock()
		fmt.Printf("Cliente %d recebeu: %+v\n", id, resp)
	}
}

// --- Funções de simulação ---
func simulateRegister(conn net.Conn, id int) {
	user := shared.User{
		UserName: fmt.Sprintf("teste%d", id),
		Password: "123",
	}
	err := utils.SendRequest(conn, "REGISTER", user)
	if err != nil {
		fmt.Printf("Cliente %d erro ao cadastrar: %v\n", id, err)
		return
	}
	fmt.Printf("Cliente %d enviou REGISTER\n", id)
}

func simulateLogin(conn net.Conn, id int) {
	user := shared.User{
		UserName: fmt.Sprintf("teste%d", id),
		Password: "123",
	}
	err := utils.SendRequest(conn, "LOGIN", user)
	if err != nil {
		fmt.Printf("Cliente %d erro ao fazer login: %v\n", id, err)
		return
	}
	fmt.Printf("Cliente %d enviou LOGIN\n", id)
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
	fmt.Printf("Cliente %d enviou PLAY\n", id)
}

func idToName(id int) string {
	return fmt.Sprintf("teste%d", id)
}
