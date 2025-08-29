package main

import (
	"fmt"
	"net"
	"PBL/server/handlers"
)

func main() {
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



