package database

import (
	"database/sql"
	"log"
	"power4/models"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() {
	datafile := "database/data.db"

	connect, err := sql.Open("sqlite3", datafile)
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

	createTables()
}

func createTables() {
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
