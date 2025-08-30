package utils

import (
	"PBL/shared"
	"bufio"
	"os"
)

var initialCards = []string{"AGUA","TERRA", "FOGO", "AR"}

func Cadastro() shared.User{
	var user shared.User
	reader := bufio.NewReader(os.Stdin)

	print("Insira o nome do usuário: ")
	user.UserName = ReadLine(reader)
	print("Insira a senha desejada: ")
	user.Password = ReadLine(reader)

	user.Cards = initialCards

	print(user.UserName)
	print(user.Password)

	return user
}

func Login() shared.User{
	var user shared.User
	reader := bufio.NewReader(os.Stdin)
	print("Insira o nome do usuário: ")
	user.UserName = ReadLine(reader)
	print("Insira a sua senha: ")
	user.Password = ReadLine(reader)
	return user
}

