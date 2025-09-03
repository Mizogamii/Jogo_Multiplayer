package services

import (
	"PBL/server/storage"
	"PBL/shared"
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

type Cliente struct {
	Connection net.Conn
	User    string
	Login     bool
	Status string
	Cards []string
}

var (
	listUsersOnline []*Cliente
	listUsersLock   sync.Mutex
	mux sync.Mutex
)

func AddUsers(cliente *Cliente) {
	listUsersLock.Lock()
	listUsersOnline = append(listUsersOnline, cliente)
	listUsersLock.Unlock()
	fmt.Println("Usuários online:", GetUsersOnline())
}

func DelUsers(cliente *Cliente) {
	listUsersLock.Lock()
	defer listUsersLock.Unlock()
	for i, j := range listUsersOnline {
		if j == cliente { // compara ponteiros
			listUsersOnline = append(listUsersOnline[:i], listUsersOnline[i+1:]...)
			break
		}
	}
	fmt.Println("Usuários online:", GetUsersOnline())
}

func CheckUser(newUser shared.User) bool {
	user, err := storage.LoadUser(newUser.UserName)
	if err != nil {
		fmt.Println("Erro ao carregar usuários:", err) 
		//Quando é com o primeiro cadastro ele acaba mostrando a mensagem de erro pois eu carrego os dados pra RAM primeiro e no caso como é o primeiro cadastro não tem dados para serem lidos. Mas como eu preciso verificar se o cadastro existe para não ter dois usuários iguais preciso fazer a leitura antes. Depois pensa em uma forma de resolver isso.
		return false
	}
	
	if user.UserName == newUser.UserName && user.Password == newUser.Password {
		return true
	}
	return false
}

func SendMessage(senderUser string, receiver string, message string) string {
	listUsersLock.Lock()
	defer listUsersLock.Unlock()

	for _, i := range listUsersOnline {
		if i.User == receiver {
			_, err := i.Connection.Write([]byte(fmt.Sprintf("%s: %s\n", senderUser, message)))
			if err != nil {
				return "ERRO: Falha ao enviar a mensagem\n"
			}
			return "OK"
		}
	}

	return "ERRO: Usuário não está online.\n"
}

func SendResponse(conn net.Conn, status string, message string, data interface{}) {
    resp := shared.Response{
        Status:  status,
        Message: message,
        Data:    data,
    }

    json.NewEncoder(conn).Encode(resp)
}

func SetStatus(username, status string){
	mux.Lock()
	defer mux.Unlock()

	for _, i := range listUsersOnline{
		if i.User == username{
			i.Status = status
			break
		}
	}
}


func GetUsersOnline() []string {
	names := []string{}
	listUsersLock.Lock()
	defer listUsersLock.Unlock()
	for _, c := range listUsersOnline {
		names = append(names, c.User)
	}
	return names
}

func GetClientByName(userName string) *Cliente{
    listUsersLock.Lock()
    defer listUsersLock.Unlock()

    for _, c := range listUsersOnline {
        if c.User == userName {
			fmt.Println("O nome do cara ai: ", c)
            return c
        }
    }
    return nil 
}

func UserOnline(userName string) bool{
	online := GetUsersOnline()
	for _, i := range online{
		if i == userName {
			return true
		}
	}
	return false
}

func GetClientByConn(conn net.Conn) *Cliente {
	listUsersLock.Lock()
	defer listUsersLock.Unlock()

	for _, c := range listUsersOnline {
		if c.Connection == conn {
			return c
		}
	}
	return nil
}



