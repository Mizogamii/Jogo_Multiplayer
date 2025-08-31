package game

import (
	"PBL/server/models"
	"PBL/server/services"
	"fmt"
	"math/rand"
)

func CreateRoom(player1, player2 *services.Cliente){
	var turn *services.Cliente
	if rand.Intn(2) == 0{
		turn = player1
	}else{
		turn = player2
	}
	
	room := &models.Room{
		Player1: player1,
		Player2: player2,
		Turn: turn,
	}

	fmt.Println(room)

	models.GameRooms[player1.User] = room
	models.GameRooms[player2.User] = room
}