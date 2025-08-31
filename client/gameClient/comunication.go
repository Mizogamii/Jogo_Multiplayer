package gameClient

import (
	"PBL/client/utils"
	"PBL/shared"
	"bufio"
	"fmt"
	"os"
)
func ShowGame(user shared.User) string{
	utils.ListCards(user)
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Insira a carta desejada: ")
	input := utils.ReadLine(reader)
	fmt.Println(input) 

	return input
}