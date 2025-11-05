package database

import (
	"database/sql"
	"log"
	"power4/models"

	_ "modernc.org/sqlite"
)

func InitDB() {
	datafile := "database/data.db"

	connect, err := sql.Open("sqlite", datafile)
	if err != nil {
		log.Fatal("Failed to connect to SQLite:", err)
	}

	if err = connect.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("SQLite Connected")

	models.DB = &models.Database{
		Connect:  connect,
		Datafile: datafile,
	}

	createUserTable()
	createHistoryTable()
}

func createUserTable() {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		email TEXT,
		password_hash TEXT NOT NULL,
		elo INTEGER DEFAULT 1000,
		victoires INTEGER DEFAULT 0,
		defaites INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := models.DB.Connect.Exec(query)
	if err != nil {
		log.Println("Failed to create table", err)
	}
}

func createHistoryTable() {
    query := `
    CREATE TABLE IF NOT EXISTS match_history (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        player1 TEXT NOT NULL,
        player2 TEXT NOT NULL,
        winner TEXT,
        delta INTEGER DEFAULT 0,
        ranked BOOLEAN DEFAULT 0,
        date DATETIME DEFAULT CURRENT_TIMESTAMP
    );`
    _, err := models.DB.Connect.Exec(query)
    if err != nil {
        log.Println("Failed to create table", err)
    }
}