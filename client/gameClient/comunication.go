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
			return "EXITROOM"
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
			utils.ListCards(currentUser)
			reader := bufio.NewReader(os.Stdin)

			//Reinicia o deck antes de montar
			currentUser.Deck = []string{}

			cardsChosen := make(map[int]bool) //para não repetir cartas no deck

			for len(currentUser.Deck) < 4 {
				fmt.Printf("Escolha a carta %d (1 a %d): ", len(currentUser.Deck)+1, len(currentUser.Cards))
				input := utils.ReadLine(reader)
				inputInt, err := strconv.Atoi(input)
				if err != nil {
					fmt.Println("Entrada inválida! Digite um número.")
					continue
				}

				cardIndex := inputInt - 1
				if cardIndex < 0 || cardIndex >= len(currentUser.Cards) {
					fmt.Println("Número inválido! Tente novamente.")
					continue
				}

				if cardsChosen[cardIndex] {
					fmt.Println("Carta já escolhida! Escolha outra.")
					continue
				}

				//Adiciona a carta no deck
				cardsChosen[cardIndex] = true
				currentUser.Deck = append(currentUser.Deck, currentUser.Cards[cardIndex])
				fmt.Println("Carta adicionada:", currentUser.Cards[cardIndex])
			}

			fmt.Println("Deck montado com sucesso:", currentUser.Deck)

		case "4":
			return
		default:
			fmt.Println("Opção inválida!")
		}
	}
}
