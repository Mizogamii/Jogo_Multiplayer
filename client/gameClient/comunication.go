package gameClient

import (
	"PBL/client/utils"
	"PBL/shared"
	"fmt"
	"net"
	"strconv"
	"sync"
)

var result string

	func StartGame(conn net.Conn, currentUser shared.User, respChan chan shared.Response) bool {
	fmt.Printf("\n%sPartida iniciada!%s\n",utils.Purple, utils.Reset)

	turnChan := make(chan bool, 1)    
	gameOver := make(chan struct{})   

	var ended bool
	var mu sync.Mutex 
	var once sync.Once 
	
	exitRequested := false

	//Função para fechar canais de forma segura
	closeChannels := func() {
		once.Do(func() {
			close(gameOver)
			close(turnChan)
		})
	}

	//Função para marcar como terminado
	setEnded := func() bool {
		mu.Lock()
		defer mu.Unlock()
		if ended {
			return true
		}
		ended = true
		return false
	}

	//Função para verificar se terminou
	isEnded := func() bool {
		mu.Lock()
		defer mu.Unlock()
		return ended
	}

	//Goroutine para ler mensagens do servidor
	go func() {
		defer closeChannels()

		for resp := range respChan {
			if isEnded() {
				return
			}

			switch resp.Status {
			case "yourTurn":
				if !isEnded() {
					fmt.Printf("\n%sSua vez: %s%s\n", utils.Blue, currentUser.UserName, utils.Reset)
					select {
					case turnChan <- true: 
					case <-gameOver:
					}
				}

			case "opponentPlayed":
				if !isEnded() {
					fmt.Printf("\n%sOponente jogou: %s%s\n", utils.Green, resp.Data, utils.Reset)
				}

			case "gameResult", "gamefinalResult", "gameResultExit":
				if !isEnded() {
					result = resp.Message
					utils.PrintResult(result)
				}

			case "gameOver":
				if setEnded() { 
					return
				}
				fmt.Println("\n", resp.Message)
				return
			}
		}
	}()

	//Loop principal de jogadas
	for {
		select {
		case <-turnChan:
			if isEnded() {
				return exitRequested
			}

			card := ShowGame(currentUser, gameOver)
		
			if card == "EXITROOM" {
				if setEnded() { 
					return exitRequested
				}
				
				utils.SendRequest(conn, "EXITROOM", "Cliente saiu da partida")
				fmt.Printf("\n%sVocê saiu da partida.%s\n",utils.Red, utils.Reset)
				fmt.Printf("Sua vez: ", )
				exitRequested = true
				closeChannels()
				return exitRequested
			}

			if !isEnded() {
				err := utils.SendRequest(conn, "CARD", card)
				if err != nil {
					fmt.Println("Erro ao enviar carta:", err)
				}
			}

		case <-gameOver:
			return exitRequested
		}
	}
}

func ShowGame(user shared.User, gameOver chan struct{}) string {
	for {
		utils.ListCardsDeck(&user)

		fmt.Print("Insira a carta desejada (0 para sair): ")
		
		inputChan := make(chan string, 1) 
		go func() {
			defer close(inputChan)
			inputChan <- utils.ReadLineSafe()
		}()

		select {
		case <-gameOver:
			return "EXITROOM"
		case input := <-inputChan:
			inputInt, err := strconv.Atoi(input)
			if err != nil {
				fmt.Println("Entrada inválida! Digite um número entre 0 e 4.")
				continue
			}
			if inputInt < 0 || inputInt > 4 {
				fmt.Println("Número inválido! Escolha entre 0 e 4.")
				continue
			}
			utils.Clear()
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
}

func ChoiceDeck(conn net.Conn, currentUser *shared.User) {
	for {
		option := utils.ShowMenuDeck()
		utils.Clear()
	
		switch option {
		case "1":
			utils.ListCards(currentUser)
		case "2":
			utils.ListCardsDeck(currentUser)
		case "3":
			utils.ListCards(currentUser)

			currentUser.Deck = []string{}
			cardsChosen := make(map[int]bool)

			for len(currentUser.Deck) < 4 {
				fmt.Printf("Escolha a carta %d (1 a %d): ", len(currentUser.Deck)+1, len(currentUser.Cards))
				input := utils.ReadLineSafe()
				inputInt, err := strconv.Atoi(input)
				if err != nil || inputInt < 1 || inputInt > len(currentUser.Cards) {
					fmt.Println("Número inválido! Tente novamente.")
					continue
				}

				cardIndex := inputInt - 1
				if cardsChosen[cardIndex] {
					fmt.Println("Carta já escolhida! Escolha outra.")
					continue
				}

				cardsChosen[cardIndex] = true
				currentUser.Deck = append(currentUser.Deck, currentUser.Cards[cardIndex])
				fmt.Println("Carta adicionada:", currentUser.Cards[cardIndex])
			}

			fmt.Println("Deck montado com sucesso:", currentUser.Deck)
			err := utils.SendRequest(conn, "DECK", *currentUser)
			if err != nil {
				fmt.Println("Erro ao salvar deck no servidor:", err)
			}

		case "4":
			utils.Clear()
			return
		default:
			fmt.Println("Opção inválida!")
		}
	}
}