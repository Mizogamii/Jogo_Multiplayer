package utils

import (
	"PBL/shared"
	"fmt"
)

func Cadastro() shared.User {
	var user shared.User

	fmt.Print("Insira o nome do usuário: ")
	user.UserName = ReadLineSafe()
	fmt.Print("Insira a senha desejada: ")
	user.Password = ReadLineSafe()

	return user
}

func Login() shared.User {
	var user shared.User
	fmt.Print("Insira o nome do usuário: ")
	user.UserName = ReadLineSafe()
	print("Insira a sua senha: ")
	user.Password = ReadLineSafe()

	return user
}
