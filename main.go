package main

import (
	"log"
	"net/http"
	"power4/database"
	"power4/models"
	"power4/routes"
)

func main() {
	// Initialize the database connection and create tables if they don't exist
	database.InitDB()
	defer models.DB.Connect.Close()

	// Register all HTTP handlers for the application
	routes.SetupRoutes()

	// Start the HTTP server and listen for incoming requests
	log.Print("Listening on localhost:8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
