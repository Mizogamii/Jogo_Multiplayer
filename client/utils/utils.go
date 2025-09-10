package utils

import (
	"PBL/shared"
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"
)

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


func ShowWaitingScreen(stopChan chan bool) {
    frames := []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}
    i := 0
    for {
        select {
        case <-stopChan:
            fmt.Println("\nPartida encontrada!")
            return
        default:
            fmt.Printf("\r%s Procurando partida%s", frames[i%len(frames)], strings.Repeat(".", i%4))
            i++
            time.Sleep(100 * time.Millisecond)
        }
    }
}

//Função para fazer input com espaçamentos e etc
func ReadLine(reader *bufio.Reader) string {
    for {
        text, err := reader.ReadString('\n')
        if err != nil {
            continue
        }
        text = strings.TrimSpace(text)
        if text != "" {
            return text
        }
    }
}

func SendRequest(conn net.Conn, action string, data interface{}) error {
    //converte data para json.RawMessage
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

func Clear(){
    fmt.Print("\033[H\033[2J")
}