package gameClient

import (
	"PBL/client/utils"
	"PBL/server/storage"
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
        case "gameOver":
            fmt.Println("Cabou")
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

func ShowGame(user shared.User) string{
	utils.ListCards(user)
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Insira a carta desejada: ")
	input := utils.ReadLine(reader)	
	
	switch input{
	case "1":
		return  "AGUA"
	case "2":
		return  "TERRA"
	case "3":
		return "FOGO"
	case "4":
		return "AR"
	}
	return ""
}

func ChoiceDeck(conn net.Conn, currentUser shared.User) {
	for {
		option := utils.ShowMenuDeck()

		switch option {
		case "1":
			//Lista todas as cartas do jogador
			utils.ListCards(currentUser)

		case "2":
			//Lista apenas as cartas do deck atual
			utils.ListCadsDeck(currentUser)

		case "3":
			//Montar novo deck
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

			err := utils.SendRequest(conn, "UPDATE_DECK", currentUser)
			if err != nil{
				fmt.Println("Erro ao atualizar dados.")
			}
		case "4":
			//Voltar ao menu principal
			return

		default:
			fmt.Println("Opção inválida!")
		}
	}
}
