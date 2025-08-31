package game

import (
	"PBL/server/models"
	"PBL/server/services"
	"fmt"
	"net"
)

func HandleRound(room *models.Room, client *services.Cliente, card string) {
	services.SendResponse(client.Connection, "cardReceived", "Recebi a carta meu parceiro", nil)
	fmt.Println("Teste debug aaaaaaa denrto do hadle rouns")
	var opponentConn net.Conn
    if room.Player1.Connection == client.Connection {
        opponentConn = room.Player2.Connection
    } else {
        opponentConn = room.Player1.Connection
    }

	services.SendResponse(opponentConn, "opponentPlayed", "Oponente jogou uma carta", card)
}