package gameClient

import (
	"PBL/client/utils"
	"PBL/shared"
	"bufio"
	"fmt"
	"os"
	"net"
)

func StartGame(conn net.Conn, currentUser shared.User, respChan chan shared.Response){
	fmt.Println("Partida iniciada!")
	turnChan := make(chan bool) // canal para sinalizar que é a vez do jogador

// goroutine que lê do servidor
go func() {
    for resp := range respChan {
        switch resp.Status {
        case "yourTurn":
            turnChan <- true // sinaliza que é a vez
        case "opponentPlayed":
            fmt.Println("Oponente jogou:", resp.Data)
        case "gameResult":
            fmt.Println("Resultado:", resp.Message)
        case "gameOver":
            fmt.Println("Cabou")
            return
        }
    }
}()

for {
    <-turnChan 
    card := ShowGame(currentUser) // aqui bloqueia só pra entrada do jogador
    err := utils.SendRequest(conn, "CARD", card)
    if err != nil {
        fmt.Println("Erro ao enviar carta:", err)
    }
}

}

func ShowGame(user shared.User) string{
	utils.ListCards(user)
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Insira a carta desejada: ")
	input := utils.ReadLine(reader)	
	
	switch input{
	case "1":
		return  "AGUA"
	case "2":
		return  "TERRA"
	case "3":
		return "FOGO"
	case "4":
		return "AR"
	}
	return ""
}