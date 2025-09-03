package game

import "fmt"

func CheckWinner(card1, card2 string) string{
	fmt.Println("Eita a lógica")
	switch card1{
	case "FOGO":
		if card2 == "AR" || card2 == "FOGO"{
			return "EMPATE"
		}else if card2 == "TERRA"{
			return "GANHOU"
		}else if card2 == "AGUA"{
			return "PERDEU"
		}
	case "AGUA":
		if card2 == "TERRA" || card2 == "AGUA"{
			return "EMPATE"
		}else if card2 == "FOGO"{
			return "GANHOU"
		}else if card2 == "AR"{
			return "PERDEU"
		}
	case "TERRA":
		if card2 == "AGUA" || card2 == "TERRA"{
			return "EMPATE"
		}else if card2 == "AR"{
			return "GANHOU"
		}else if card2 == "FOGO"{
			return "PERDEU"
		}
	case "AR":
		if card2 == "FOGO" || card2 == "AR"{
			return "EMPATE"
		}else if card2 == "AGUA"{
			return "GANHOU"
		}else if card2 == "TERRA"{
			return "PERDEU"
		}
	}
	return ""
}