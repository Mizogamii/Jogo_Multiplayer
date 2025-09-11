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
		Mu:    sync.Mutex{},
	}
)

func HandleConnection(conn net.Conn) {
	defer conn.Close()
	decoder := json.NewDecoder(conn)

	for {
		var req shared.Request
		err := decoder.Decode(&req)

		if err != nil {
			client := services.GetClientByConn(conn)
			if client != nil {
				if err == io.EOF {
					fmt.Println("Conexão fechada de:", client.User)
				} else {
					fmt.Println("Conexão perdida de:", client.User, "-", err)
				}
				services.DelUsers(client)
			} else {
				if err == io.EOF {
					fmt.Println("Conexão fechada de cliente desconhecido")
				} else {
					fmt.Println("Erro na conexão de cliente desconhecido: ", err)
				}
			}
			return
		}

		fmt.Println("Usuários online:", services.GetUsersOnline())

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
			HandlePack(conn, req)

		case "CARD":
			var card string
			client := services.GetClientByConn(conn)

			if client == nil {
				fmt.Println("Cliente não encontrado")
				break
			}

			room := models.GameRooms[client.User]
			if room == nil {
				fmt.Println("Sala não encotrada para o cliente: ", client.User)
				break
			}
			dataBytes, _ := json.Marshal(req.Data)
			json.Unmarshal(dataBytes, &card)

			game.HandleRound(room, client, card)

		case "LOGOUT":
			fmt.Println("Logout pedido")
			client := services.GetClientByConn(conn)
			if client != nil {
				client.Login = false
				client.Status = "livre"
				services.DelUsers(client)
				services.SendResponse(conn, "successLogout", "Você foi deslogado.", nil)
			}
			return
		case "EXIT":
			fmt.Println("Saindo...")
			client := services.GetClientByConn(conn)
			if client != nil {
				client.Login = false
				client.Status = "livre"
				services.DelUsers(client)
				services.SendResponse(conn, "successExit", "Você saiu.", nil)
			}
			return

		case "EXITROOM":
			fmt.Println("Fim do jogo por desistencia...")
			client := services.GetClientByConn(conn)
			if client == nil {
				fmt.Println("EXITROOM: cliente não encontrado para conexão")
				break
			}

			room := models.GameRooms[client.User]
			if room == nil {
				fmt.Println("EXITROOM: sala não encontrada para o cliente:", client.User)
				services.SendResponse(conn, "successExitRoom", "Você saiu da sala.", nil)
				break
			}

			var opponent *services.Cliente
			if room.Player1 == client {
				opponent = room.Player2
			} else {
				opponent = room.Player1
			}

			room.Status = models.Finished
			
			services.SendResponse(client.Connection, "gameResultExit", "Você saiu da partida", nil)
			if opponent != nil {
				services.SendResponse(opponent.Connection, "gameResultExit", "Oponente desistiu — você venceu!", nil)
				services.SendResponse(opponent.Connection, "gameOver", "Fim de jogo, voltando ao menu...", nil)
				opponent.Status = "livre"
			}

			client.Status = "livre"

			delete(models.GameRooms, room.Player1.User)
			delete(models.GameRooms, room.Player2.User)
		
		case "PING":
			fmt.Println("PING")
			services.SendResponse(conn, "PONG", "Retornando", nil)

		case "LEAVEQUEUE":
			HandleLeave(conn, req)

		default:
			fmt.Println("Ação desconhecida recebida:", req.Action)
			return
		}
	}
}

func HandleRegister(conn net.Conn, req shared.Request) {
    var newUser shared.User

    data, _ := json.Marshal(req.Data)
    json.Unmarshal(data, &newUser)

    fmt.Println("☻Usuário recebido: ", newUser.UserName)

    newUser.Cards = []string{"AGUA", "TERRA", "FOGO", "AR", "MATO"}
    newUser.Deck  = []string{"AGUA", "TERRA", "FOGO", "AR"}

    exist := services.UserExist(newUser)
    if !exist {
        err := storage.SaveUsers(newUser)
        if err != nil {
            services.SendResponse(conn, "error", "Falha ao salvar usuário.", nil)
        } else {
            services.SendResponse(conn, "successRegister", "Cadastro realizado com sucesso", newUser)
        }
    } else {
        services.SendResponse(conn, "error", "Usuário já existe", nil)
    }
}

func HandleLogin(conn net.Conn, req shared.Request) (*services.Cliente, bool) {
	var user shared.User

	data, _ := json.Marshal(req.Data)
	json.Unmarshal(data, &user)

	exists := services.CheckLogin(user)
	if exists {
		if !services.UserOnline(user.UserName) {
			cliente := &services.Cliente{
				Connection: conn,
				User:       user.UserName,
				Login:      true,
				Status:     "livre",
				Password:   user.Password,
				Cards:      loadCards(user.UserName, conn),
				Deck:       loadDeck(user.UserName, conn),
			}
			services.AddUsers(cliente)
			fmt.Println(cliente.Status)

			services.SendResponse(conn, "successLogin", "Login realizado com sucesso.", shared.User{
				UserName: cliente.User,
				Cards:    cliente.Cards,
				Deck:     cliente.Deck,
			})
			return cliente, true
		}
	}
	services.SendResponse(conn, "error", "Login ou senha inválidos.", nil)
	return nil, false
}

func HandlePlay(conn net.Conn, req shared.Request) {
	var user shared.User

	if err := json.Unmarshal(req.Data, &user); err != nil {
		fmt.Println("Erro ao desserializar User:", err)
		services.SendResponse(conn, "error", "Falha ao ler dados do usuário", nil)
		return
	}

	userName := user.UserName
	fmt.Println("Nome do usuário:", userName)

	client := services.GetClientByName(userName)
	if client == nil {
		fmt.Println("Cliente não logado:", userName)
		services.SendResponse(conn, "error", "Usuário não está logado", nil)
		return
	}

	fmt.Println("Status cliente:", client.Status)

	Matchmaking.Mu.Lock()
	defer Matchmaking.Mu.Unlock()

	if client.Status == "livre" {
		client.Status = "fila"
		Matchmaking.Queue = append(Matchmaking.Queue, client)
		fmt.Println("Cliente entrou na fila:", client.User)
		services.SendResponse(conn, "successPlay", "Você entrou na fila de jogo", client)
	} else {
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

func HandleDeck(conn net.Conn, req shared.Request) {
	var user shared.User

	dataBytes, _ := json.Marshal(req.Data)
	if err := json.Unmarshal(dataBytes, &user); err != nil {
		services.SendResponse(conn, "error", "Erro ao processar dados do deck.", nil)
		return
	}

	fmt.Println("Atualizando deck do usuário:", user.UserName)
	
	client := services.GetClientByName(user.UserName)
    
    client.Deck = user.Deck

	userToSave := shared.User{
		UserName: client.User,
		Password: client.Password,
		Cards:    client.Cards,
		Deck:     client.Deck,
	}

	err := storage.SaveUsers(userToSave)

	if err != nil {
		services.SendResponse(conn, "error", "Falha ao salvar deck.", nil)
		return
	}

	services.SendResponse(conn, "successUpdateDeck", "Deck atualizado com sucesso!", userToSave)
}

func HandlePack(conn net.Conn, req shared.Request) {
	client := services.GetClientByConn(conn)
	if client == nil {
		fmt.Println("Cliente não encontrado para abrir pacote")
		return
	}

	cards, err := game.OpenPack(client.User)
	if err != nil {
		services.SendResponse(conn, "error", err.Error(), nil)
		return
	}

	client.Cards = append(client.Cards, cards...)

	userToSave := shared.User{
		UserName: client.User,
		Password: client.Password,
		Cards:    client.Cards,
		Deck:     client.Deck,
	}

	err = storage.SaveUsers(userToSave)
	if err != nil {
		services.SendResponse(conn, "error", "Falha ao salvar pack.", nil)
		return
	}
	
	services.SendResponse(conn, "successPack", "Pacote aberto com sucesso!", userToSave)
}

func HandleLeave(conn net.Conn, req shared.Request) {
    client := services.GetClientByConn(conn)
    if client == nil {
        fmt.Println("Cliente não encontrado para sair da fila")
        services.SendResponse(conn, "error", "Cliente não encontrado", nil)
        return
    }

    Matchmaking.Mu.Lock()
    defer Matchmaking.Mu.Unlock()

    found := false
    for i, c := range Matchmaking.Queue {
        if c == client {
            Matchmaking.Queue = append(Matchmaking.Queue[:i], Matchmaking.Queue[i+1:]...)
            client.Status = "livre"
            found = true
            break
        }
    }

    if found {
        fmt.Println("Cliente saiu da fila:", client.User)
        services.SendResponse(conn, "successLeaveQueue", "Você saiu da fila", nil)
    } else {
        services.SendResponse(conn, "error", "Você não estava na fila", nil)
    }

    names := []string{}
    for _, c := range Matchmaking.Queue {
        names = append(names, c.User)
    }
    fmt.Println("Fila atual:", names)
}

//Cria partidas entre jogadores na fila (matchmaking).
func StartMatchmaking() {
	for {
		Matchmaking.Mu.Lock()

		if len(Matchmaking.Queue) >= 2 {
			player1 := Matchmaking.Queue[0]
			player2 := Matchmaking.Queue[1]

			fmt.Printf("DEBUG Matchmaking - Encontrando match: %s vs %s\n", player1.User, player2.User)

			Matchmaking.Queue = Matchmaking.Queue[2:]

			player1.Status = "jogando"
			player2.Status = "jogando"

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

//Notificação que conseguiu um oponente para ambos os lados
func notifyClient(player1, player2 *services.Cliente) {
	services.SendResponse(player1.Connection, "match", "Oponente encontrado", player2.User)
	services.SendResponse(player2.Connection, "match", "Oponente encontrado", player1.User)
}
//Carregar a lista de cartas que o cliente possui
func loadCards(userName string, conn net.Conn) []string {
	user, err := storage.LoadUser(userName)
	if err != nil {
		fmt.Println("Erro ao carregar usuários:", err)
		return nil
	}
	return user.Cards
}

//Carregar a lista de cartas no deck do cliente
func loadDeck(userName string, conn net.Conn) []string {
	user, err := storage.LoadUser(userName)
	if err != nil {
		fmt.Println("Erro ao carregar usuários:", err)
		return nil
	}

	return user.Deck
}
