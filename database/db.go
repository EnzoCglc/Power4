package database

import (
	"database/sql"
	"log"
	"power4/models"

	_ "modernc.org/sqlite"
)

// InitDB initializes the SQLite database connection and creates required tables.
func InitDB() {
	datafile := "database/data.db"

	// Open connection to SQLite database
	connect, err := sql.Open("sqlite", datafile)
	if err != nil {
		log.Fatal("Failed to connect to SQLite:", err)
	}

	// Verify the connection is working
	if err = connect.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("SQLite Connected")

	// Store connection in global variable for use throughout the application
	models.DB = &models.Database{
		Connect:  connect,
		Datafile: datafile,
	}

	// Create tables if they don't exist
	createUserTable()
	createHistoryTable()
}

// createUserTable creates the users table if it doesn't already exist.
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

// createHistoryTable creates the match_history table if it doesn't already exist.
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