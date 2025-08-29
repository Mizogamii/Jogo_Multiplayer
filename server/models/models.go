package models

import(
	"PBL/server/services"
    "sync"
)

type Matchmaking struct {
    Queue []*services.Cliente
    Mu sync.Mutex
}

type Room struct {
    Player1 *services.Cliente
    Player2 *services.Cliente
}