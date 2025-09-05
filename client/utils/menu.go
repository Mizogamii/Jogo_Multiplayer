package utils

import (
	"PBL/shared"
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func Menu(conn net.Conn) string{
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\n--------------------------------")
	fmt.Println("           Menu Inicial           ")
	fmt.Println("--------------------------------")
	fmt.Println("1 - Cadastro")
	fmt.Println("2 - Login")
	fmt.Println("3 - Sair")
	fmt.Print("Insira a opção desejada: ")
	option := ReadLine(reader)

	switch option {
	case "1": 
		return "REGISTER"

	case "2":
		return "LOGIN"

	case "3":
		fmt.Println("Saindo...") 
		conn.Close() 
		return "EXIT"
		
	default:
		fmt.Print("ERRO! Digite apenas números de 1 a 3")
		return ""
	}
}

func ShowMenuLogin(conn net.Conn) string{
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\n--------------------------------")
		fmt.Println("              Menu              ")
		fmt.Println("--------------------------------")
		fmt.Println("1 - Entrar na fila")
		fmt.Println("2 - Ver/alterar deck")
		fmt.Println("3 - Abrir pacote")
		fmt.Println("4 - Deslogar")
		fmt.Print("Insira a opção desejada: ")
		input := ReadLine(reader)
		return input
	}
}

func ShowMenuDeck() string{
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\n--------------------------------")
		fmt.Println("            Menu deck           ")
		fmt.Println("--------------------------------")
		fmt.Println("1 - Visualizar todas as cartas")
		fmt.Println("2 - Visualizar cartas do deck")
		fmt.Println("3 - Alterar o deck")
		fmt.Println("4 - Voltar ao menu principal")
		fmt.Print("Insira a opção desejada: ")
		input := ReadLine(reader)
		return input
}


func ListCards(user shared.User) {
	fmt.Println("\n--------------------------------")
	fmt.Println("             Cartas             ")
	fmt.Println("--------------------------------")
	for i, card := range user.Cards {
		fmt.Printf("%d: %s\n", i+1, card)
	}
	fmt.Println("--------------------------------")
}

func ListCardsDeck(user shared.User) {
	fmt.Println("\n--------------------------------")
	fmt.Println("             Deck             ")
	fmt.Println("--------------------------------")
	for i, card := range user.Deck {
		fmt.Printf("%d: %s\n", i+1, card)
	}
	fmt.Println("--------------------------------")
}



//AS FUNÇÕES DAQUI PRA BAIXO DEVERIAM IR PARA OUTRO CANTO, ESSE AQUI É SÓ PARA MENUS
func ListenServer(conn net.Conn, respChan chan shared.Response, stopChan chan bool) {
    decoder := json.NewDecoder(conn)
    for {
        var resp shared.Response
        if err := decoder.Decode(&resp); err != nil {
            fmt.Println("Erro ao receber mensagem do servidor:", err)
            close(respChan)
            return
        }
        
        respChan <- resp
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
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func SendRequest(conn net.Conn, action string, data interface{}) error {
    // converte data para json.RawMessage
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