package main

import (
	"PBL/client/gameClient"
	"PBL/client/utils"
	"PBL/shared"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	//Para docker
	//serverAddress := "servidor:8080"

	//Para teste local
	serverAddress := "localhost:8080" 

	conn, err := net.Dial("tcp", serverAddress) 

	if err != nil {
		fmt.Println("Erro ao conectar:", err)
		return
	}
	defer conn.Close()

	inputMgr := utils.GetInputManager()
    inputMgr.Start()
    defer inputMgr.Stop()

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
				utils.Clear()
			
			case "LOGIN":
				requestData = utils.Login()
				utils.Clear()

			case "EXIT":
				fmt.Printf("\n%sSaindo...%s\n", utils.Cyan,utils.Reset)
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
			utils.Clear()
			switch operationTypeLogin {
			//Entrar na fila e jogar
			case "1":
				err = utils.SendRequest(conn, "PLAY", currentUser)
				if err != nil {
					fmt.Println("Erro:", err)
					continue
				}
				
				stopChan := make(chan bool)
				go utils.ShowWaitingScreen(conn, stopChan)
				
				queue := true
				for queue{
					resp := <-respChan
					/*fmt.Printf("DEBUG - Resposta completa: Status='%s', Message='%s', Data='%v', Tipo Data: %T\n", resp.Status, resp.Message, resp.Data, resp.Data)*/
	
					//Deu match --> mostra a tela da partida
					if resp.Status == "match"{
						queue = false
						stopChan <- true
						exitRequested := gameClient.StartGame(conn, currentUser, respChan)
						if exitRequested {
								fmt.Println("\nVoltando ao menu principal...")
						}
							break // Sai do loop de aguardar match
					}else if resp.Status == "successLeaveQueue"{
						queue = false
						fmt.Println("Você saiu da fila. Voltando ao menu...")
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
					fmt.Printf("Resposta antiga descartada: %s\n", oldResp.Status)
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
				
			case "5":
				utils.ShowPing()
				start := time.Now()
				
				utils.SendRequest(conn, "PING", "Mandando ping")

				<-respChan //recebendo a resposta mas ignora para só pegar o tempo

				elapsed := time.Since(start)

    			fmt.Println("Tempo de ping:", elapsed.Nanoseconds(), "ns")
				fmt.Println("\n----------------------------------")

			//Deslogar
			case "6":
				fmt.Printf("\n%sDeslogado com sucesso!%s\n", utils.Cyan, utils.Reset)
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