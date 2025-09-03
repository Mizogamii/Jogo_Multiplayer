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
	data, err := json.MarshalIndent(newUser, "", "  ")
	if err != nil {
		return err
	}

	dir := GetDataDir()

	//Cria a pasta se não existir
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