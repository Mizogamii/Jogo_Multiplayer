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
	}else{
		opponent = room.Player1
	}
	
	services.SendResponse(opponent.Connection, "opponentPlayed", "Oponente jogou uma carta", card)

	room.Turn = opponent
	services.SendResponse(opponent.Connection, "yourTurn", "Sua vez parceiro ☻", nil)
}