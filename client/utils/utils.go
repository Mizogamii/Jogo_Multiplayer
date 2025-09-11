package utils

import (
	"PBL/shared"
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

type InputManager struct {
	scanner    *bufio.Scanner
	inputChan  chan string
	stopChan   chan struct{}
	mu         sync.Mutex
	isRunning  bool
}

var (
	inputManager     *InputManager
	inputManagerOnce sync.Once
)

func GetInputManager() *InputManager {
	inputManagerOnce.Do(func() {
		inputManager = &InputManager{
			scanner:   bufio.NewScanner(os.Stdin),
			inputChan: make(chan string, 10), 
			stopChan:  make(chan struct{}),
		}
	})
	return inputManager
}

func (im *InputManager) Start() {
	im.mu.Lock()
	defer im.mu.Unlock()

	if im.isRunning {
		return
	}
	im.isRunning = true
	im.stopChan = make(chan struct{}) 
	go im.readInput()
}

func (im *InputManager) Stop() {
	im.mu.Lock()
	defer im.mu.Unlock()

	if !im.isRunning {
		return
	}
	close(im.stopChan) 
	im.isRunning = false
}

func (im *InputManager) readInput() {
	for {
		select {
		case <-im.stopChan:
			return
		default:
			if im.scanner.Scan() {
				text := strings.TrimSpace(im.scanner.Text())
				if text != "" {
					select {
					case im.inputChan <- text: 
					default:
						<-im.inputChan
						im.inputChan <- text
					}
				}
			} else if err := im.scanner.Err(); err != nil {
			
				im.scanner = bufio.NewScanner(os.Stdin)
			}
		}
	}
}

func (im *InputManager) ReadLine() string {
	return <-im.inputChan
}

func (im *InputManager) ReadLineWithTimeout(timeout time.Duration) (string, bool) {
	select {
	case text := <-im.inputChan:
		return text, true
	case <-time.After(timeout):
		return "", false
	}
}


func ListenServer(conn net.Conn, respChan chan shared.Response, stopChan chan bool) {
	decoder := json.NewDecoder(conn)
	for {
		select {
		case <-stopChan:
			fmt.Println("Encerrando ListenServer")
			return
		default:
			var resp shared.Response
			if err := decoder.Decode(&resp); err != nil {
				fmt.Println("Erro ao receber mensagem do servidor:", err)
				close(respChan)
				return
			}
			respChan <- resp
		}
	}
}

func ShowWaitingScreen(conn net.Conn, stopChan chan bool) {
	frames := []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}
	i := 0

	inputMgr := GetInputManager()

	for {
		select {
		case <-stopChan:
			fmt.Println("\nPartida encontrada!")
			return

		default:
			if text, ok := inputMgr.ReadLineWithTimeout(150 * time.Millisecond); ok {
				if text == "0" {
					fmt.Println("\nSaindo da fila...")
					SendRequest(conn, "LEAVEQUEUE", nil)
					return
				}
			}

			fmt.Printf("\r%s Procurando partida%s", frames[i%len(frames)], strings.Repeat(".", i%4))
			i++
		}
	}
}


func SendRequest(conn net.Conn, action string, data interface{}) error {
	rawData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("erro ao converter data para JSON: %w", err)
	}

	req := shared.Request{
		Action: action,
		Data:   rawData,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("erro ao converter request para JSON: %w", err)
	}

	_, err = conn.Write(jsonData)
	if err != nil {
		return fmt.Errorf("erro ao enviar para o servidor: %w", err)
	}

	return nil
}

func Clear() {
	nameOS := runtime.GOOS
	fmt.Println("Sistema operacional:", nameOS)

	var cmd *exec.Cmd
	if nameOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func Pause() {
	fmt.Print("Pressione ENTER para continuar...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func ReadLineSafe() string {
	return GetInputManager().ReadLine()
}