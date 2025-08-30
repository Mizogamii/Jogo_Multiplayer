package main

import (
	"PBL/client/utils"
	"PBL/shared"
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

	//Canais de comuniação
	stopChan := make(chan bool)
	respChan := make(chan shared.Response)

	go utils.ListenServer(conn, respChan, stopChan)

	var loginOk bool
	var user interface{}	

	for{
		if !loginOk{
			operationType := utils.Menu()

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

			err := utils.SendRequest(conn, operationType, user) 
			if err != nil{
				fmt.Println("Erro:", err)
				continue
			}

			resp := <-respChan 

			fmt.Println("Resposta do servidor:", resp.Status, resp.Message)

			if resp.Status == "successLogin"{
				loginOk = true
			}else{
				fmt.Println("ERRO Login inválido: ", resp.Message)
			}

		}else{
			operationTypeLogin := utils.ShowMenuLogin(conn)
			loginOk = true
			var action string

			switch operationTypeLogin{
			case "1":
				action = "PLAY"
				fmt.Println("PLAY")
				go utils.ShowWaitingScreen(stopChan)

			case "2":
				utils.ListCards(user)
				action = "DECK"
				fmt.Println("DECK")

			case "3":
				action = "PACK"
				fmt.Println("PACK")

			case "4":
				action = "EXIT"
				fmt.Println("Deslogado com sucesso!")
				os.Exit(1)

			default:
				fmt.Println("ERRO: Digite apenas números de 1 a 4")
				continue
			}

			err := utils.SendRequest(conn, action, user) 
			if err != nil{
				fmt.Println("Erro:", err)
				continue
			}

			resp := <-respChan

			if resp.Status == "match" {
				fmt.Printf("Oponente encontrado: %v\n", resp.Data)
			} else {
				fmt.Println("Resposta do servidor:", resp.Status, resp.Message, resp.Data)
			}
		}
	}	
}
