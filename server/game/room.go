package game

import (
	"PBL/server/models"
	"PBL/server/services"
	"fmt"
	"math/rand"
)

func CreateRoom(player1, player2 *services.Cliente) *models.Room {
	var turn *services.Cliente
	if rand.Intn(2) == 0 {
		turn = player1
	} else {
		turn = player2
	}

	room := &models.Room{
		Player1: player1,
		Player2: player2,
		Turn:    turn,
		Status: models.InProgress,
	}

	fmt.Println(room)

	models.GameRooms[player1.User] = room
	models.GameRooms[player2.User] = room

	services.SendResponse(player1.Connection, "match", "Oponente encontrado", player2.User)
	services.SendResponse(player2.Connection, "match", "Oponente encontrado", player1.User)
	fmt.Println("Enviando match")

	services.SendResponse(turn.Connection, "yourTurn", "É sua vez de jogar!", nil)
	fmt.Println("Enviando yourTurn")

	return room
}
