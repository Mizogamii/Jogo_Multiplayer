package utils

import (
	"bufio"
	"os"
)


type User struct{
	UserName string
	Password string 
}

func Cadastro() User{
	var user User
	reader := bufio.NewReader(os.Stdin)

	print("Insira o nome do usuário: ")
	user.UserName = ReadLine(reader)
	print("Insira a senha desejada: ")
	user.Password = ReadLine(reader)

	print(user.UserName)
	print(user.Password)

	return user
}

func Login() User{
	var user User
	reader := bufio.NewReader(os.Stdin)
	print("Insira o nome do usuário: ")
	user.UserName = ReadLine(reader)
	print("Insira a sua senha: ")
	user.Password = ReadLine(reader)
	return user
}

