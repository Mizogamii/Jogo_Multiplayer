package main

import (
	"PBL/server/game"
	"PBL/server/handlers"
	"fmt"
	"net"
)

func main() {

	game.BuildGlobalDeck() //Contruindo deck global para os pacotes

	go handlers.StartMatchmaking()

	listener, err := net.Listen("tcp", ":8080")

	if err != nil {
		fmt.Println("Erro ao iniciar servidor:", err)
		return
	}

	fmt.Println("Servidor rodando na porta 8080...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Erro ao aceitar a conexão: ", err)
			continue
		}
		go handlers.HandleConnection(conn)
	}
}



