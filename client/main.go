package main

import (
	"PBL/client/utils"
	"encoding/json"
	"fmt"
	"net"
	"PBL/shared"
)

func main() {
	//conn, err := net.Dial("tcp", "servidor:8080") //para docker
	conn, err := net.Dial("tcp", "localhost:8080") //para teste local
	if err != nil {
		fmt.Println("Erro ao conectar:", err)
		return
	}
	defer conn.Close()
	
	for{
		operationType := utils.Menu()
		var user interface{}

		switch operationType{
		case "REGISTER":
			user = utils.Cadastro()
			fmt.Println("Usuário cadastrado: ", user)

		case "LOGIN":
			user = utils.Login()
			fmt.Println("Usuário logado: ", user)
		
		case "EXIT":
			fmt.Println("Saindo...")
			return
		default:
			fmt.Println("ERRO: Opção inválida.")
		}

		request := shared.Request{
			Action: operationType,
			Data: user,
		}

		jsonData, err := json.Marshal(request)
		if err != nil{
			fmt.Println("Erro ao converter para json: ", err)
			return
		}

		_, err = conn.Write(jsonData)
		if err != nil{
			fmt.Println("Erro ao enviar para o servidor: ", err)
			return
		}

		var resp shared.Response
    	decoder := json.NewDecoder(conn)
    	err = decoder.Decode(&resp)
    	if err != nil {
        	fmt.Println("Erro ao ler resposta:", err)
        	return
    	}

    	fmt.Println("Resposta do servidor:", resp.Status, resp.Message)

		if resp.Status == "success"{
			utils.ShowMenuLogin(conn)
		}else{
			fmt.Println("ERRO Login inválido: ", resp.Message)
		}
	}

}
