package game

import (
	"PBL/server/models"
	"PBL/server/services"
	"fmt"
)

func CreateRoom(player1, player2 *services.Cliente){
	room := models.Room{
		Player1: player1,
		Player2: player2,
	}
	fmt.Println(room)
}