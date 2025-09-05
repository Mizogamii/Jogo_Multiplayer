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
			room.ScoreP1++
			services.SendResponse(room.Player1.Connection, "gameResult", "Ganhou!☻", nil)
			services.SendResponse(room.Player2.Connection, "gameResult", "Perdeu!", nil)
		case "PERDEU":
			//PLAYER 1 PERDEU
			//PLAYER 2 GANHOU
			room.ScoreP2++
			services.SendResponse(room.Player1.Connection, "gameResult", "Perdeu!", nil)
			services.SendResponse(room.Player2.Connection, "gameResult", "Ganhou!☻", nil)
		case "P1-EXIT":
			//PLAYER 1 DESISTIU
			//PLAYER 2 GANHOU
			room.Rounds = 3 //Para encerrar o jogo
			services.SendResponse(room.Player1.Connection, "gameResultExit", "Desistiu", nil)
			services.SendResponse(room.Player2.Connection, "gameResultExit", "Ganhou!☻", nil)
			return
			
		case"P2-EXIT":
			//PLAYER 1 GANHOU
			//PLAYER 2 DESISTIU
			room.Rounds = 3 //Para encerrar já que desistiu
			services.SendResponse(room.Player1.Connection, "gameResultExit", "Ganhou!☻", nil)
			services.SendResponse(room.Player2.Connection, "gameResultExit", "Desistiu", nil)
			return
		}
		room.Rounds++
		room.CardP1 = ""
		room.CardP2 = ""

		if room.Rounds >= 3 || room.ScoreP1 == 2 || room.ScoreP2 == 2 {
			if room.ScoreP1 > room.ScoreP2 {
				services.SendResponse(room.Player1.Connection, "gamefinalResult", "Vitória final!", nil)
				services.SendResponse(room.Player2.Connection, "gamefinalResult", "Derrota final!", nil)

			} else if room.ScoreP2 > room.ScoreP1 {
				services.SendResponse(room.Player1.Connection, "gamefinalResult", "Derrota final!", nil)
				services.SendResponse(room.Player2.Connection, "gamefinalResult", "Vitória final!", nil)

			} else {
				services.SendResponse(room.Player1.Connection, "gamefinalResult", "Empate final!", nil)
				services.SendResponse(room.Player2.Connection, "gamefinalResult", "Empate final!", nil)
			}
		}
	}
}
