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
	for{
		resp := <- respChan
		switch resp.Status{
		case "yourTurn":
			fmt.Println("Sua vez! Pode jogar meu parceiro!")
			
			card := ShowGame(currentUser)
			err := utils.SendRequest(conn, "CARD", card)
			fmt.Println("Mandou a carta: ", card)
			if err != nil {
				fmt.Println("Erro ao enviar carta:", err)
				continue	
			}
		case "opponentPlayed":
			fmt.Println("Oponente jogou: ", resp.Data)
		
		case "gameOver":
			fmt.Println("Cabou")
			return
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