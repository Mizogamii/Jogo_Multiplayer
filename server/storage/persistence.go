package storage

import (
	"PBL/shared"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func getDataDir() string {
	cwd, _ := os.Getwd()
	return filepath.Join(cwd, "data")
}

func getUsersFilePath() string {
	return filepath.Join(getDataDir(), "users.json")
}

func SaveUsers(newUser shared.User) error {
	users, err := LoadUsers()
	if err != nil {
		return err
	}

	users = append(users, newUser)

	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}

	dir := getDataDir()
	//Cria a pasta se não existir
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("falha ao criar pasta: %w", err)
	}

	filePath := getUsersFilePath()
	fmt.Println("Salvando em:", filePath)

	return os.WriteFile(filePath, data, 0644)
}

func LoadUsers() ([]shared.User, error) {
	users := []shared.User{}
	filePath := getUsersFilePath()

	//Pasta existe
	err := os.MkdirAll(getDataDir(), os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("falha ao criar pasta: %w", err)
	}

	//Lê o arquivo JSON
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return users, nil
		}
		return nil, err
	}

	if len(data) == 0 {
		fmt.Println("Arquivo vazio")
		return users, nil
	}

	err = json.Unmarshal(data, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}


