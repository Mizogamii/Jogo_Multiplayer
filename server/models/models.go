package models

import(
	"PBL/server/services"
    "sync"
)

var GameRooms = make(map[string]*Room)

type RoomStatus string

const (
    InProgress RoomStatus = "INPROGRESS"
    Finished RoomStatus = "FINISHED"
)

type Matchmaking struct {
    Queue []*services.Cliente
    Mu sync.Mutex
}

type Room struct {
    Player1 *services.Cliente
    Player2 *services.Cliente
    Status RoomStatus
    Turn *services.Cliente
    CardP1 string
    CardP2 string

    ScoreP1 int
    ScoreP2 int

    Rounds int
}

