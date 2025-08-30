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
		/*Status: ,
		Turn: */
	}

	fmt.Println(room)
	player1.Status = "jogando"
	player2.Status = "jogando"
}