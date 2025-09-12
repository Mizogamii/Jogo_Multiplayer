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
	Password string
	Cards []string
	Deck []string
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
	for i, j := range listUsersOnline {
		if j == cliente { 
			listUsersOnline = append(listUsersOnline[:i], listUsersOnline[i+1:]...)
			break
		}
	}
	listUsersLock.Unlock()
	fmt.Println("Usuários online:", GetUsersOnline())
}


func CheckLogin(newUser shared.User) bool {
	user, err := storage.LoadUser(newUser.UserName)
	if err != nil {
		fmt.Println("Erro ao carregar usuários:", err) 
		return false
	}
	
	if user.UserName == newUser.UserName && user.Password == newUser.Password {
		return true
	}
	return false
}

func UserExist(newUser shared.User) bool {
    user, err := storage.LoadUser(newUser.UserName)
    if err != nil {
        if err.Error() == "user not found" {
            return false
        }
        fmt.Println("Erro ao carregar usuário:", err)
        return false
    }
    return user.UserName == newUser.UserName
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


