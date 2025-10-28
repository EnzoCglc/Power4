package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3" 
)

var DB *sql.DB
const datafile string = "database/data.db" 

func InitDB(){
	var err error

	DB,err = sql.Open("sqlite3", datafile)
	if err != nil {
		log.Println("Failed to connect Sqlite", err)
	}

	if err = DB.Ping(); err != nil {
		log.Println("Failed to ping db ")
	}

	log.Println("Sqlite Connected")

	createTable()
}

func createTable() {
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

	_ , err := DB.Exec(query)
	if err != nil {
		log.Println("Failed to create table", err)
	}
}