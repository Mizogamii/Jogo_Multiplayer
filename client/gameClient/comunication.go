package gameClient

import (
	"PBL/client/utils"
	"PBL/shared"
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
)

func StartGame(conn net.Conn, currentUser shared.User, respChan chan shared.Response) bool {
	fmt.Println("Partida iniciada!")

	turnChan := make(chan bool, 1)    // sinal de vez do jogador (buffered)
	gameOver := make(chan struct{})   // sinal de fim de jogo

	var ended bool
	var mu sync.Mutex // protege a variável ended
	var once sync.Once // garante que os canais sejam fechados apenas uma vez

	exitRequested := false

	// Função para fechar canais de forma segura
	closeChannels := func() {
		once.Do(func() {
			close(gameOver)
			close(turnChan)
		})
	}

	// Função para marcar como terminado
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
					select {
					case turnChan <- true: // avisa que é a vez
					case <-gameOver: // se já acabou, não trava
					}
				}

			case "opponentPlayed":
				if !isEnded() {
					fmt.Println("Oponente jogou:", resp.Data)
				}

			case "gameResult", "gamefinalResult", "gameResultExit":
				if !isEnded() {
					fmt.Println("Resultado:", resp.Message)
				}

			case "gameOver":
				if setEnded() { // se já estava ended, não processa
					return
				}
				fmt.Println("Fim de jogo, voltando ao menu...")
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
				fmt.Println("Você saiu da partida.")
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
	reader := bufio.NewReader(os.Stdin)

	for {
		utils.ListCardsDeck(&user)
		fmt.Print("Insira a carta desejada (0 para sair): ")

		inputChan := make(chan string, 1) 
		go func() {
			defer close(inputChan)
			inputChan <- utils.ReadLine(reader)
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

		switch option {
		case "1":
			utils.ListCards(currentUser)
		case "2":
			utils.ListCardsDeck(currentUser)
		case "3":
			utils.ListCards(currentUser)
			reader := bufio.NewReader(os.Stdin)
			currentUser.Deck = []string{}
			cardsChosen := make(map[int]bool)

			for len(currentUser.Deck) < 4 {
				fmt.Printf("Escolha a carta %d (1 a %d): ", len(currentUser.Deck)+1, len(currentUser.Cards))
				input := utils.ReadLine(reader)
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
			return
		default:
			fmt.Println("Opção inválida!")
		}
	}
}
