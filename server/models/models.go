package models

import (
	"net"
)

type User struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type Cliente struct {
	Connection net.Conn
	User    string
	Login     bool
	Status string
}

type Response struct {
	Status string
	Message string `json:"message,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username"`
    Password string `json:"password"`
}

type PlayRequest struct {
	Move string `json:"move"`
}

type ChatMessage struct {
	From string `json:"from"`
	To string `json:"to"`
	Content string `json:"content"`
}