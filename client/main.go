package main

import (
	"PBL/client/gameClient"
	"PBL/client/utils"
	"PBL/shared"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

func main() {
	//conn, err := net.Dial("tcp", "servidor:8080") //para docker
	//Modifiquei demais e agora nao tá ridando pelo docker CONSERTE DEPOIS
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
	var currentUser shared.User

	for {
		if !loginOk {
			operationType := utils.Menu()
			var requestData interface{}

			switch operationType {
			case "REGISTER":
				requestData = utils.Cadastro()

			case "LOGIN":
				requestData = utils.Login()

			case "EXIT":
				fmt.Println("Saindo...")
				return
			default:
				fmt.Println("ERRO: Opção inválida.")
				continue
			}

			err = utils.SendRequest(conn, operationType, requestData)
			if err != nil {
				fmt.Println("Erro:", err)
				continue
			}

			resp := <-respChan

			fmt.Println("\n", resp.Message)

			if resp.Status == "successLogin" {
				loginOk = true

				var serverUser shared.User

				dataBytes, _ := json.Marshal(resp.Data)
				if err := json.Unmarshal(dataBytes, &serverUser); err != nil {
					fmt.Println("Erro ao desserializar dados do login: ", err)
					continue
				}

				currentUser = shared.User{
					UserName: serverUser.UserName,
					Cards:    serverUser.Cards,
					Deck:     []string{},
				}

			} else if resp.Status != "successRegister" && resp.Status != "successLogin"{
				fmt.Println("ERRO Login inválido: ", resp.Message)
			}

		} else {
			operationTypeLogin := utils.ShowMenuLogin(conn)
			var action string

			switch operationTypeLogin {
			case "1":
				action = "PLAY"
				fmt.Println("PLAY")
				go utils.ShowWaitingScreen(stopChan)
				err = utils.SendRequest(conn, action, currentUser)
				if err != nil {
					fmt.Println("Erro:", err)
					continue
				}
				for {
					resp := <-respChan
					fmt.Printf("DEBUG - Resposta completa: Status='%s', Message='%s', Data='%v', Tipo Data: %T\n", resp.Status, resp.Message, resp.Data, resp.Data)

					//Deu match --> mostra a tela da partida
					switch resp.Status{
					case "match":
						stopChan <- true
						gameClient.StartGame(conn, currentUser, respChan)
					case "successPlay":
						fmt.Println("Aguardando oponente...")
					case "opponentPlayed":
						fmt.Println("Oponente jogou:", resp.Data)
					default:
						fmt.Println("Resposta inesperada:", resp.Status)
					}
				}

			case "2": //tá com problema RESOLVE
				gameClient.ChoiceDeck(currentUser)
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

			if err != nil {
				fmt.Println("Erro ao converter currentUser para JSON:", err)
				continue
			}
			if action != "PLAY" {
				err = utils.SendRequest(conn, action, currentUser)

				if err != nil {
					fmt.Println("Erro:", err)
					continue
				}
			}

			resp := <-respChan

			fmt.Println("Printando resp: ", resp.Status)

		}
	}
}
