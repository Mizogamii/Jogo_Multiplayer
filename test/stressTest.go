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

type Metrics struct {
	Latencies  []time.Duration
	QueueTimes []time.Duration
	Messages   int
	mu         sync.Mutex
}

var metrics = Metrics{}

var (
	CreateRoomsTest = make(map[string]TestRoom) // agora é um map
	QueueTest       []string
	mu              sync.Mutex
	wg              sync.WaitGroup
)

func main() {
	numClients := 100
	wg.Add(numClients)

	for i := 1; i <= numClients; i++ {
		go simulateClient(i)
	}

	wg.Wait()

	// Relatório final
	fmt.Println("----- RELATÓRIO FINAL -----")
	fmt.Println("Salas criadas:")
	i := 1
	for _, sala := range CreateRoomsTest {
		fmt.Printf("Sala %d: %v\n", i, sala.Players)
		i++
	}
	fmt.Println("Clientes ainda na fila:", QueueTest)
	fmt.Println("---------------------------")

	// Relatório de desempenho
	reportMetrics()
}

func simulateClient(id int) {
	defer wg.Done()

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Printf("Cliente %d não conseguiu conectar: %v\n", id, err)
		return
	}
	defer conn.Close()

	start := time.Now()

	// canal para receber respostas do servidor
	respChan := make(chan shared.Response, 10)
	stopChan := make(chan bool)
	go utils.ListenServer(conn, respChan, stopChan)

	// cadastro
	simulateRegister(conn, id)
	time.Sleep(50 * time.Millisecond)

	// login
	simulateLogin(conn, id)
	time.Sleep(50 * time.Millisecond)

	// entra na fila
	simulateEnterQueue(conn, id)

	timeout := time.After(5 * time.Second) // cliente escuta por no máximo 5s

	for {
		select {
		case resp, ok := <-respChan:
			if !ok {
				return
			}

			metrics.mu.Lock()
			metrics.Messages++
			metrics.mu.Unlock()

			switch resp.Status {
			case "match":
				elapsed := time.Since(start)
				metrics.mu.Lock()
				metrics.Latencies = append(metrics.Latencies, elapsed)
				metrics.mu.Unlock()

				p1 := idToName(id)
				p2 := fmt.Sprintf("%v", resp.Data)

				// ordena os nomes para evitar duplicação
				if p1 > p2 {
					p1, p2 = p2, p1
				}

				key := fmt.Sprintf("%s|%s", p1, p2)

				mu.Lock()
				if _, ok := CreateRoomsTest[key]; !ok {
					CreateRoomsTest[key] = TestRoom{Players: []string{p1, p2}}
				}
				mu.Unlock()

				fmt.Printf("Cliente %d entrou em partida: %+v\n", id, resp)
				return

			case "successPlay":
				mu.Lock()
				QueueTest = append(QueueTest, idToName(id))
				mu.Unlock()
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

// --- Relatório de métricas ---
func reportMetrics() {
	metrics.mu.Lock()
	defer metrics.mu.Unlock()

	var total time.Duration
	for _, l := range metrics.Latencies {
		total += l
	}

	avgLatency := time.Duration(0)
	if len(metrics.Latencies) > 0 {
		avgLatency = total / time.Duration(len(metrics.Latencies))
	}

	fmt.Println("===== RELATÓRIO DE DESEMPENHO =====")
	fmt.Printf("Mensagens recebidas: %d\n", metrics.Messages)
	fmt.Printf("Latência média: %v\n", avgLatency)
	fmt.Printf("Latência mínima: %v\n", min(metrics.Latencies))
	fmt.Printf("Latência máxima: %v\n", max(metrics.Latencies))
	fmt.Println("===================================")
}

func min(vals []time.Duration) time.Duration {
	if len(vals) == 0 {
		return 0
	}
	m := vals[0]
	for _, v := range vals {
		if v < m {
			m = v
		}
	}
	return m
}

func max(vals []time.Duration) time.Duration {
	if len(vals) == 0 {
		return 0
	}
	m := vals[0]
	for _, v := range vals {
		if v > m {
			m = v
		}
	}
	return m
}
