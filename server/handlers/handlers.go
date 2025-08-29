package handlers

import (
	"PBL/server/services"
	"PBL/server/storage"
	"PBL/shared"
	"encoding/json"
	"fmt"
	"net"
	"io"
)

func HandleConnection(conn net.Conn) {
	defer conn.Close()
	decoder := json.NewDecoder(conn)

	for {
		var req shared.Request
		err := decoder.Decode(&req)

		if err != nil {
			if err == io.EOF {
        		fmt.Println("Cliente desconectou")
        		return
			}
			fmt.Println("Erro ao ler ou decodificar JSON:", err)
			return
		}

		switch req.Action  {
		case "REGISTER":
    		HandleRegister(conn, req)

		case "LOGIN":
    		HandleLogin(conn, req)
    
		default:
			return
		}
	}
}

func HandleRegister(conn net.Conn, req shared.Request) {
	var user shared.User

	data, _ := json.Marshal(req.Data)
	json.Unmarshal(data, &user)

	fmt.Println("Usuário recebido: ", user.UserName)

	exists := services.CheckUser(user)
	if !exists {
		err := storage.SaveUsers(user)
		if err != nil {
			services.SendResponse(conn, "error", "Falha ao salvar usuário.", nil)
	
		} else {
			services.SendResponse(conn, "success", "Cadastro realizado", nil)
			fmt.Println("Cadastro ok")
		}
	} else {
		services.SendResponse(conn, "error", "Usuário já existe", nil)
	}
}

func HandleLogin(conn net.Conn, req shared.Request) (*services.Cliente, bool) {
	var user shared.User

	data, _ := json.Marshal(req.Data)
	json.Unmarshal(data, &user)

	exists := services.CheckUser(user)
	if exists {
		fmt.Println("Existe ai")
		if !services.UserOnline(user.UserName) {
			cliente := &services.Cliente{
				Connection: conn,
				User:       user.UserName,
				Login:      true,
				Status:     "livre",
			}
			services.AddUsers(cliente)
			fmt.Println("Login ok")
			services.SendResponse(conn, "success", "Login realizado com sucesso.", nil)
			return cliente, true
		}
	}
	services.SendResponse(conn, "error", "Login ou senha inválidos.", nil)
	return nil, false
}
