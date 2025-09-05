package storage

import (
	"PBL/shared"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func GetDataDir() string {
	cwd, _ := os.Getwd()
	return filepath.Join(cwd, "data")
}

func getUsersFilePath(userName string) string {
	return filepath.Join(GetDataDir(), userName+".json")
}

func SaveUsers(newUser shared.User) error {
	// Carregar usuário existente
	oldUser, err := LoadUser(newUser.UserName)
	if err == nil {
		//Se o campo de senha veio vazio, mantém a antiga
		if newUser.Password == "" {
			newUser.Password = oldUser.Password
		}
		//Se o deck não veio, mantém o antigo
		if len(newUser.Deck) == 0 {
			newUser.Deck = oldUser.Deck
		}
	}

	data, err := json.MarshalIndent(newUser, "", "  ")
	if err != nil {
		return err
	}

	dir := GetDataDir()

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("falha ao criar pasta: %w", err)
	}

	filePath := getUsersFilePath(newUser.UserName)
	fmt.Println("Salvando em:", filePath)

	return os.WriteFile(filePath, data, 0644)
}

func LoadUser(userName string) (shared.User, error) {
	var user shared.User
	filePath := getUsersFilePath(userName)

	data, err := os.ReadFile(filePath)
	if err != nil{
		return user, err
	}

	if len(data) == 0{
		fmt.Println("Arquivo vazio")
		return user, nil
	}

	err = json.Unmarshal(data, &user)
	if err != nil {
		return user, nil
	}
	return user, nil
}