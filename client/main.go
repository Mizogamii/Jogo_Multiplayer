package main

import (
	"PBL/client/utils"
	"fmt"
	"net"
	"os"
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

		case "LOGIN":
			user = utils.Login()
		
		case "EXIT":
			fmt.Println("Saindo...")
			return
		default:
			fmt.Println("ERRO: Opção inválida.")
		}

		resp, err := utils.SendRequest(conn, operationType, user) 
		if err != nil{
			fmt.Println("Erro:", err)
    		continue
		}

    	fmt.Println("Resposta do servidor:", resp.Status, resp.Message)

		if resp.Status == "success"{
			operationTypeLogin := utils.ShowMenuLogin(conn)
			var action string
			switch operationTypeLogin{
			case "1":
				action = "PLAY"
				fmt.Println("PLAY")

			case "2":
				action = "PACK"
				fmt.Println("PACK")

			case "3":
				action = "EXIT"
				fmt.Println("Deslogado com sucesso!")
				os.Exit(1)

			default:
				fmt.Println("ERRO: Digite apenas números de 1 a 3")
				continue
			}

			resp, err := utils.SendRequest(conn, action, user) 
			if err != nil{
				fmt.Println("Erro:", err)
				continue
			}
			fmt.Println("Resposta do servidor:", resp.Status, resp.Message)

	}else{
			fmt.Println("ERRO Login inválido: ", resp.Message)
		}
	}

}
