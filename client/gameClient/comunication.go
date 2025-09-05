package gameClient

import (
	"PBL/client/utils"
	"PBL/shared"
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
)

func StartGame(conn net.Conn, currentUser shared.User, respChan chan shared.Response){
	fmt.Println("Partida iniciada!")
	turnChan := make(chan bool) 

	go func() {
		for resp := range respChan {
			switch resp.Status {
			case "yourTurn":
				turnChan <- true 

			case "opponentPlayed":
				fmt.Println("Oponente jogou:", resp.Data)

			case "gameResult":
				fmt.Println("Resultado:", resp.Message)

			case "gameResultExit": 
				fmt.Println("Resultado:", resp.Message)
				return

			case "gamefinalResult":
				fmt.Println("Fim do jogo!")
				fmt.Println("Resultado:", resp.Message)
				return
			}
		}
	}()

	for {
		<-turnChan 
		card := ShowGame(currentUser) 
		err := utils.SendRequest(conn, "CARD", card)
		if err != nil {
			fmt.Println("Erro ao enviar carta:", err)
		}
	}
}

func ShowGame(user shared.User) string {
	reader := bufio.NewReader(os.Stdin)

	for {
		utils.ListCardsDeck(user)
		fmt.Print("Insira a carta desejada (0 para sair): ")

		input := utils.ReadLine(reader)
		inputInt, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Entrada inválida! Digite um número entre 0 e 4.")
			continue
		}

		if inputInt < 0 || inputInt > 4 {
			fmt.Println("Número inválido! Escolha entre 0 e 4.")
			continue
		}

		switch inputInt {
		case 0:
			return "EXIT"
		case 1:
			return user.Deck[0]
		case 2:
			return user.Deck[1]
		case 3:
			return user.Deck[2]
		case 4:
			return user.Deck[3]
		}
	}
}


func ChoiceDeck(currentUser shared.User) {
	for {
		option := utils.ShowMenuDeck()

		switch option {
		case "1":
			//Lista todas as cartas do jogador
			utils.ListCards(currentUser)

		case "2":
			//Lista apenas as cartas do deck atual
			utils.ListCardsDeck(currentUser)

		case "3":
			//Montar novo deck
			utils.ListCardsDeck(currentUser)
			reader := bufio.NewReader(os.Stdin)

			//Reinicia o deck antes de montar
			currentUser.Deck = []string{}

			for i := 0; i < 4; i++ {
				fmt.Print("Insira o número da carta desejada: ")
				input := utils.ReadLine(reader)
				inputInt, err := strconv.Atoi(input)
				if err != nil {
					fmt.Println("Erro ao converter para número:", err)
					i-- 
					continue
				}

				cardIndex := inputInt - 1

				//Validação do índice
				if cardIndex < 0 || cardIndex >= len(currentUser.Cards) {
					fmt.Println("Número inválido! Tente novamente.")
					i--
					continue
				}

				fmt.Println("Carta escolhida:", currentUser.Cards[cardIndex])
				currentUser.Deck = append(currentUser.Deck, currentUser.Cards[cardIndex])
			}

			fmt.Println("Seu deck foi montado com sucesso:", currentUser.Deck)

		case "4":
			//Voltar ao menu principal
			return

		default:
			fmt.Println("Opção inválida!")
		}
	}
}
