package models

import (
	"encoding/json"
	"os"
)

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	Elo          int    `json:"elo"`
	Win          int    `json:"win"`
	Losses       int    `json:"losses"`
}

type DataBase struct {
	Users []User `json:"users"`
}

func LoadDB(path string) (DataBase, error) {
	var db DataBase
	data, err := os.ReadFile(path)
	if err != nil {
		return db, err
	}

	err = json.Unmarshal(data, &db)
	return db, err
}

func SaveDB(path string, db DataBase) error {
	data, err := json.MarshalIndent(db, "", "")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
