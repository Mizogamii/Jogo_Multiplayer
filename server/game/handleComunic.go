package game

import (
	"PBL/server/models"
	"PBL/server/services"
	"fmt"
)

func HandleRound(room *models.Room, client *services.Cliente, card string) {
	if room.Turn != client{
		services.SendResponse(client.Connection, "error", "Não é sua vez", nil)
		return
	}

	fmt.Println(client.User + "jogou a carta: "+ card)

	var opponent *services.Cliente

	if room.Player1 == client{
		opponent = room.Player2
		room.CardP1 = card
		fmt.Println("Carta jogada P1: ", room.CardP1)
		
	}else{
		opponent = room.Player1
		room.CardP2 = card
		fmt.Println("Carta jogada P2: ", room.CardP2)
	}
	
	services.SendResponse(opponent.Connection, "opponentPlayed", "Oponente jogou uma carta", card)

	room.Turn = opponent

	services.SendResponse(opponent.Connection, "yourTurn", "Sua vez parceiro ☻", nil)
	
	if room.CardP1 != "" && room.CardP2 != "" {
		result := CheckWinner(room.CardP1, room.CardP2)

		switch result {
		case "EMPATE":
			services.SendResponse(room.Player1.Connection, "gameResult", "Empate!", nil)
			services.SendResponse(room.Player2.Connection, "gameResult", "Empate!", nil)
		case "GANHOU":
			//PLAYER 1 GANHOU
			//PLAYER 2 PERDEU
			services.SendResponse(room.Player1.Connection, "gameResult", "Ganhou!☻", nil)
			services.SendResponse(room.Player2.Connection, "gameResult", "Perdeu!", nil)
		case "PERDEU":
			//PLAYER 1 PERDEU
			//PLAYER 2 GANHOU
			services.SendResponse(room.Player1.Connection, "gameResult", "Perdeu!", nil)
			services.SendResponse(room.Player2.Connection, "gameResult", "Ganhou!☻", nil)
		}
		room.CardP1 = ""
		room.CardP2 = ""
	}
}


