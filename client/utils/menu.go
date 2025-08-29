package utils

import (
	"PBL/shared"
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

func Menu() string{
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\n--------------------------")
	fmt.Println("       Menu Inicial       ")
	fmt.Println("--------------------------")
	fmt.Println("1 - Cadastro")
	fmt.Println("2 - Login")
	fmt.Println("3 - Sair")
	fmt.Print("Insira a opção desejada: ")
	option := ReadLine(reader)

	switch option {
	case "1": 
		return "REGISTER"
	case "2":
		return "LOGIN"
	case "3":
		fmt.Println("Saindo...")
		os.Exit(1) //isso pra não voltar pra o menu princiapl --> o usuario sai do servidor
	default:
		fmt.Print("ERRO! Digite apenas números de 1 a 3")
		return ""
	}
	return ""
}

func ShowMenuLogin(conn net.Conn){
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("--------------------------")
		fmt.Println("           Menu           ")
		fmt.Println("--------------------------")
		fmt.Println("1 - Entrar na fila")
		fmt.Println("2 - Abrir pacote")
		fmt.Println("3 - Deslogar")
		fmt.Print("Insira a opção desejada: ")
		input := ReadLine(reader)

		var action string
		switch input{
		case "1":
			action = "PLAY"
		case "2":
			action = "PACK"
		case "3":
			action = "EXIT"
			fmt.Println("Deslogado com sucesso!")
			os.Exit(1)
		default:
			fmt.Println("ERRO: Digite apenas números de 1 a 3")
			continue
		}
		
		req := shared.Request{
			Action: action,
			Data: nil,
		}

		jsonData, err := json.Marshal(req)
		if err != nil{
			fmt.Println("Erro ao converter json: ", err)	
			continue
		}

		_, err = conn.Write(jsonData)
		if err != nil{
			fmt.Println("Erro ao envar para o servidor: ", err)
			return
		}
	}
}

func ShowMenuGame(){
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("--------------------------")
	fmt.Println("           Menu           ")
	fmt.Println("--------------------------")
	fmt.Println("1 - Escolher carta") //a ideia seria aqui ser a parte de montar o deck 
	//tenho que ver que jogo vou fazer agora já que não tenho tanto tempo mais...
	//deixo o usuário escolher 5 cartas para o deck e faço lutar com esses? 
	//as cartas especiais seriam as cartas só decoradas?

	fmt.Print("Insira a opção desejada: ")
	input := ReadLine(reader)
	fmt.Println(input) //só pra não dar erro mesmo depois eu contnuo
}

//Função para fazer input com espaçamentos e etc
func ReadLine(reader *bufio.Reader) string {
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}