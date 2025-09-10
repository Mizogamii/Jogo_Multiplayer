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
			operationType := utils.Menu(conn)
			var requestData interface{}

			switch operationType {
			case "REGISTER":
				requestData = utils.Cadastro()

			case "LOGIN":
				requestData = utils.Login()

			case "EXIT":
				fmt.Println("Saindo...")		
				conn.Close()
				os.Exit(0)		
				
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
					Password: serverUser.Password,
					Cards:    serverUser.Cards,
					Deck:     serverUser.Deck,
				}

			} else if resp.Status != "successRegister" && resp.Status != "successLogin"{
				fmt.Println("ERRO: ", resp.Message)
			}

		} else {
			operationTypeLogin := utils.ShowMenuLogin(conn)

			switch operationTypeLogin {
			//Entrar na fila e jogar
			case "1":
				err = utils.SendRequest(conn, "PLAY", currentUser)
				if err != nil {
					fmt.Println("Erro:", err)
					continue
				}
	
				go utils.ShowWaitingScreen(stopChan)
				
				for{
					resp := <-respChan
					fmt.Printf("DEBUG - Resposta completa: Status='%s', Message='%s', Data='%v', Tipo Data: %T\n", resp.Status, resp.Message, resp.Data, resp.Data)
	
					//Deu match --> mostra a tela da partida
					if resp.Status == "match"{
						stopChan <- true
						exitRequested := gameClient.StartGame(conn, currentUser, respChan)
						if exitRequested {
								loginOk = false
								fmt.Println("Voltando ao menu principal...")
						} else {
							fmt.Println("Partida finalizada. Voltando ao menu...")
						}
							break // Sai do loop de aguardar match
						}
				}
				
			//Alterar deck ou só ver
			case "2": 
				gameClient.ChoiceDeck(conn, &currentUser)
				utils.SendRequest(conn, "DECK", currentUser)

			//Abrir pacote
			case "3":
				//Limpa respostas pendentes no canal 
				select {
				case oldResp := <-respChan:
					fmt.Printf("DEBUG - Resposta antiga descartada: %s\n", oldResp.Status)
				default:
					
				}
	
				err = utils.SendRequest(conn, "PACK", currentUser)
				if err != nil {
					fmt.Println("Erro ao enviar requisição PACK:", err)
					continue
				}

				resp := <-respChan
				fmt.Println("Status da resposta:", resp.Status)
				fmt.Println("Mensagem:", resp.Message)
				
				if resp.Status == "successPack" {
					var updatedUser shared.User
					dataBytes, _ := json.Marshal(resp.Data)
					if err := json.Unmarshal(dataBytes, &updatedUser); err != nil {
						fmt.Println("Erro ao desserializar dados do pacote:", err)
						continue
					}
					
					//Atualiza o currentUser com as novas cartas
					currentUser.Cards = updatedUser.Cards
					
					fmt.Println("Pacote aberto com sucesso! Novas cartas adicionadas.")
				} else {
					fmt.Println("Erro ao abrir pacote:", resp.Message)
				}
				
				utils.ListCards(&currentUser)

			//Ver as regras do jogo
			case "4":
				utils.ShowRules()
				
			//Deslogar
			case "5":
				fmt.Println("Deslogado com sucesso!")
				err = utils.SendRequest(conn, "LOGOUT", currentUser)
				if err != nil {
					fmt.Println("Erro:", err)
				}
				conn.Close()
				os.Exit(0)

			default:
				fmt.Println("ERRO: Digite apenas números de 1 a 5")
				continue
			}

			if err != nil {
				fmt.Println("Erro ao converter currentUser para JSON:", err)
				continue
			}

		}
	}
}
