package main

import (
	"log"
	"net/http"
	"power4/database"
	"power4/routes"
)

func main() {
	// Init Database
	database.InitDB()
	
	// Setup routes
	routes.SetupRoutes()

	_, err := database.DB.Exec(`INSERT INTO users (username, email) VALUES (?, ?)`, "enzo", "enzo@example.com")
	if err != nil {
		log.Println("ERROR when insert data:", err)
	}

	table, err := database.DB.Query(`SELECT id, username, email FROM users`)
	if err != nil {
		log.Println("Error when listen table")
	}

	for table.Next() {
		var id int
		var username string
		var email string
		table.Scan(&id, &username, &email)
		log.Printf("Valeur de l'id : %d  Valeur de l'username : %s Valeur de l'email : %s",id, username, email)
	}

	defer table.Close()
	defer database.DB.Close()

	// Start Server on localhost:8080
	log.Print("Listening on localhost:8080...")
	error := http.ListenAndServe(":8080", nil)
	if error != nil {
		log.Fatal(error)
	}
}
