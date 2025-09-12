package main

import (
	"PBL/client/utils"
	"PBL/shared"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"
)


func listenResponses(t *testing.T, conn net.Conn) chan shared.Response {
	respChan := make(chan shared.Response, 10)
	stopChan := make(chan bool)
	go utils.ListenServer(conn, respChan, stopChan)
	return respChan
}

func connectTestClient(t *testing.T) net.Conn {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		t.Fatalf("falha ao conectar ao servidor: %v", err)
	}
	return conn
}

//Teste do cadastro
func TestRegisterUser(t *testing.T) {
	conn := connectTestClient(t)
	defer conn.Close()

	user := shared.User{
		UserName: "testeUnitario",
		Password: "123",
	}

	err := utils.SendRequest(conn, "REGISTER", user)
	if err != nil {
		t.Errorf("erro ao cadastrar usuário: %v", err)
	} else {
		t.Log("Cadastro realizado com sucesso")
	}
}

//Teste do login
func TestLoginUser(t *testing.T) {
	conn := connectTestClient(t)
	defer conn.Close()

	user := shared.User{
		UserName: "testeUnitario",
		Password: "123",
	}

	err := utils.SendRequest(conn, "LOGIN", user)
	if err != nil {
		t.Errorf("erro ao fazer login: %v", err)
	} else {
		t.Log("Login realizado com sucesso")
	}
}

func TestQueueAndMatchmaking(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		clientFlow(t, 1)
	}()
	go func() {
		defer wg.Done()
		clientFlow(t, 2)
	}()

	wg.Wait()
}

func clientFlow(t *testing.T, id int) {
	conn := connectTestClient(t)
	defer conn.Close()

	respChan := listenResponses(t, conn)

	// Cadastro
	user := shared.User{
		UserName: fmt.Sprintf("teste%d", id),
		Password: "123",
	}
	err := utils.SendRequest(conn, "REGISTER", user)
	if err != nil {
		t.Logf("Cliente %d: usuário já existe ou erro no cadastro", id)
	}

	// Login
	err = utils.SendRequest(conn, "LOGIN", user)
	if err != nil {
		t.Errorf("Cliente %d: erro ao fazer login: %v", id, err)
		return
	}

	err = utils.SendRequest(conn, "PLAY", user)
	if err != nil {
		t.Errorf("Cliente %d: erro ao entrar na fila: %v", id, err)
		return
	}

	timeout := time.After(5 * time.Second)
	for {
		select {
		case resp := <-respChan:
			if resp.Status == "match" {
				t.Logf("Cliente %d entrou em partida com: %v", id, resp.Data)
				return
			}
		case <-timeout:
			t.Errorf("Cliente %d: timeout esperando matchmaking", id)
			return
		}
	}
}