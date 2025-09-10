package utils

import (
	"PBL/shared"
	"bufio"
	"fmt"
	"net"
	"os"
)

func Menu(conn net.Conn) string{
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\n----------------------------------")
	fmt.Println("            Menu Inicial            ")
	fmt.Println("----------------------------------")
	fmt.Println("1 - Cadastro")
	fmt.Println("2 - Login")
	fmt.Println("3 - Sair")
	fmt.Print("Insira a opção desejada: ")
	option := ReadLine(reader)
	//fmt.Println("DEBUG - input lido:", option)

	switch option {
	case "1": 
		return "REGISTER"

	case "2":
		return "LOGIN"

	case "3":
		fmt.Println("Saindo...") 
		conn.Close() 
		return "EXIT"
		
	default:
		fmt.Print("ERRO! Digite apenas números de 1 a 3")
		return ""
	}
}

func ShowMenuLogin(conn net.Conn) string{
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\n----------------------------------")
		fmt.Println("               Menu               ")
		fmt.Println("----------------------------------")
		fmt.Println("1 - Entrar na fila")
		fmt.Println("2 - Ver/alterar deck")
		fmt.Println("3 - Abrir pacote")
		fmt.Println("4 - Visualizar regras")
		fmt.Println("5 - Visualizar ping")
		fmt.Println("6 - Deslogar")
		fmt.Print("Insira a opção desejada: ")
		input := ReadLine(reader)
		//fmt.Println("DEBUG - input lido:", input)
		return input
	}
}

func ShowMenuDeck() string{
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\n--------------------------------")
		fmt.Println("            Menu deck           ")
		fmt.Println("--------------------------------")
		fmt.Println("1 - Visualizar todas as cartas")
		fmt.Println("2 - Visualizar cartas do deck")
		fmt.Println("3 - Alterar o deck")
		fmt.Println("4 - Voltar ao menu principal")
		fmt.Print("Insira a opção desejada: ")
		input := ReadLine(reader)
		return input
}


func ListCards(user *shared.User) {
	fmt.Println("\n----------------------------------")
	fmt.Println("            Suas cartas           ")
	fmt.Println("----------------------------------")
	for i, card := range user.Cards {
		fmt.Printf("%d: %s\n", i+1, card)
	}
	fmt.Println("----------------------------------")
}

func ListCardsDeck(user *shared.User) {
	fmt.Println("\n----------------------------------")
	fmt.Println("             Seu deck             ")
	fmt.Println("----------------------------------")
	for i, card := range user.Deck {
		fmt.Printf("%d: %s\n", i+1, card)
	}
	fmt.Println("----------------------------------")
}


func ShowRules(){
	fmt.Println("\n----------------------------------")
	fmt.Println("              Regras              ")
	fmt.Println("----------------------------------")
	fmt.Println("Ao fazer o cadastro você recebeu\n5 cartas. Sendo elas: AGUA, TERRA,\nFOGO, AR e MATO")
	fmt.Println("\nCada carta tem seus pontos fortes\ne fracos:")
	fmt.Println("\n ÁGUA")
	fmt.Println(" Forte contra FOGO")
	fmt.Println(" Fraco contra AR")

	fmt.Println("\n TERRA")
	fmt.Println(" Forte contra AR")
	fmt.Println(" Fraco contra FOGO")

	fmt.Println("\n FOGO")
	fmt.Println(" Forte contra TERRA")
	fmt.Println(" Fraco contra ÁGUA")

	fmt.Println("\n AR")
	fmt.Println(" Forte contra ÁGUA")
	fmt.Println(" Fraco contra TERRA")

	fmt.Println("\n MATO")
	fmt.Println(" Carta MISTERIOSA")
	
	fmt.Println("----------------------------------")


}
/*Com emoji
func ShowRules2(){
	fmt.Println("\n----------------------------------")
	fmt.Println("              Regras              ")
	fmt.Println("----------------------------------")
	fmt.Println("Ao fazer o cadastro você recebeu\n5 cartas. Sendo elas: AGUA, TERRA,\nFOGO, AR e MATO")
	fmt.Println("\nCada carta tem seus pontos fortes\ne fracos:")
	fmt.Println("\n ÁGUA💧")
	fmt.Println(" Forte contra FOGO")
	fmt.Println(" Fraco contra AR")

	fmt.Println("\n TERRA🌱")
	fmt.Println(" Forte contra AR")
	fmt.Println(" Fraco contra FOGO")

	fmt.Println("\n FOGO🔥")
	fmt.Println(" Forte contra TERRA")
	fmt.Println(" Fraco contra ÁGUA")

	fmt.Println("\n AR💨")
	fmt.Println(" Forte contra ÁGUA")
	fmt.Println(" Fraco contra TERRA")

	fmt.Println("\n MATO🌿")
	fmt.Println(" Carta MISTERIOSA")

}*/