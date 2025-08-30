package handlers

import (
	"PBL/server/game"
	"PBL/server/models"
	"PBL/server/services"
	"PBL/server/storage"
	"PBL/shared"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

var (
	Matchmaking = &models.Matchmaking{
		Queue: make([]*services.Cliente, 0),
		Mu: sync.Mutex{},
	}
)

func HandleConnection(conn net.Conn) {
	defer conn.Close()
	decoder := json.NewDecoder(conn)

	for {
		var req shared.Request
		err := decoder.Decode(&req)

		if err != nil {
			if err == io.EOF {
        		fmt.Println("Cliente desconectou")
        		return
			}
			fmt.Println("Erro ao ler ou decodificar JSON:", err)
			return
		}

		switch req.Action {
		case "REGISTER":
    		HandleRegister(conn, req)

		case "LOGIN":
    		HandleLogin(conn, req)
    
		case "PLAY":
			HandlePlay(conn, req)

		case "DECK":
			HandleDeck(conn, req)

		case "PACK":
			HandlePack()

		default:
			return
		}
	}
}

func HandleRegister(conn net.Conn, req shared.Request) {
	var user shared.User

	data, _ := json.Marshal(req.Data)
	json.Unmarshal(data, &user)

	fmt.Println("Usuário recebido: ", user.UserName)
	fmt.Println("Cartas: ", user.Cards)

	exists := services.CheckUser(user)
	if !exists {
		err := storage.SaveUsers(user)
		if err != nil {
			services.SendResponse(conn, "error", "Falha ao salvar usuário.", nil)
	
		} else {
			services.SendResponse(conn, "successRegister", "Cadastro realizado", nil)
			fmt.Println("Cadastro ok")
		}
	} else {
		services.SendResponse(conn, "error", "Usuário já existe", nil)
	}
}

func HandleLogin(conn net.Conn, req shared.Request) (*services.Cliente, bool) {
	var user shared.User
	
	data, _ := json.Marshal(req.Data)
	json.Unmarshal(data, &user)

	exists := services.CheckUser(user)
	if exists {
		fmt.Println("Existe ai")
		if !services.UserOnline(user.UserName) {
			cliente := &services.Cliente{
				Connection: conn,
				User:       user.UserName,
				Login:      true,
				Status:     "livre",
				Cards: loadCards(user.UserName, conn),
			}
			services.AddUsers(cliente)
			fmt.Println(cliente.Status)
			fmt.Println("Login ok")
			services.SendResponse(conn, "successLogin", "Login realizado com sucesso.", cliente)
			return cliente, true
		}
	}
	services.SendResponse(conn, "error", "Login ou senha inválidos.", nil)
	return nil, false
}

func HandlePlay(conn net.Conn, req shared.Request) {
	fmt.Println("Play do server uau")

	var user shared.User
	data, _ := json.Marshal(req.Data)
	json.Unmarshal(data, &user)
	userName := user.UserName

	client := services.GetClientByName(userName)
	if client == nil {
		fmt.Println("Cliente não logado: ", userName)
		services.SendResponse(conn, "error", "Usuário não está logado", nil)
		return
	}

	Matchmaking.Mu.Lock()
	defer Matchmaking.Mu.Unlock()

	if client.Status == "livre"{
		client.Status = "fila"
		Matchmaking.Queue = append(Matchmaking.Queue, client)
		fmt.Println("Cliente entrou na fila:", client.User)
		services.SendResponse(conn, "successPlay", "Você entrou na fila de jogo", nil)
	}else{
		services.SendResponse(conn, "error", "Você já está na fila", nil)
		return
	}

	//Mostra a fila atual
	names := []string{}
	for _, c := range Matchmaking.Queue {
		names = append(names, c.User)
	}
	fmt.Println("Fila atual:", names)

}

//listo todas as cartas que o usuario tem com os indices
//faço ele digitar o nome/numero da carta que ele quer no deck (5 cartas)

func HandleDeck(conn net.Conn, req shared.Request) {
	fmt.Println("deck")
}
	

func HandlePack(){
	fmt.Println("Pack do server uau")
	//tem que fazer uma função que adiociona as cartas dos pacotes no json do usuario. lembra de ffazer 

}

/*
func getUsersQueue() []string {
	names := []string{}
	Matchmaking.Mu.Lock()
	defer Matchmaking.Mu.Unlock()
	for _, c := range Matchmaking.Queue {
		names = append(names, c.User)
	}
	return names
}

func inQueue(userName string) bool{
	in := getUsersQueue()
	for _, i := range in{
		if i == userName {
			return true
		}
	}
	return false
}*/


func StartMatchmaking(){
	for{
		Matchmaking.Mu.Lock()

		if len(Matchmaking.Queue) >= 2{
			player1 := Matchmaking.Queue[0]
			player2 := Matchmaking.Queue[1]

			Matchmaking.Queue = Matchmaking.Queue[2:]
			
			Matchmaking.Mu.Unlock()

			fmt.Println("Queue: ", Matchmaking.Queue)

			go notifyClient(player1, player2)
			go game.CreateRoom(player1, player2)
			continue
		}
		Matchmaking.Mu.Unlock()
		time.Sleep(100 * time.Millisecond)
	}
}

//LEMBRA DE ATUALIZAR O STATUS DO CLIENTE QUANDO ACABAR O JOGO
func notifyClient(player1, player2 *services.Cliente){
	go services.SendResponse(player1.Connection, "match", "Oponente encontrado", player2.User)
	go services.SendResponse(player2.Connection, "match", "Oponente encontrado", player1.User)
}

func loadCards(userName string, conn net.Conn,) []string{
	user, err := storage.LoadUser(userName)
	if err != nil {
		fmt.Println("Erro ao carregar usuários:", err)
		return nil
	}
	
	return user.Cards
}