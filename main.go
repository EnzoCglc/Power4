package main

import (
	"log"
	"net/http"
	"power4/database"
	"power4/models"
	"power4/routes"
)

func main() {
	// Init Database
	database.InitDB()
	defer models.DB.Connect.Close()

	// Setup routes
	routes.SetupRoutes()

	// Start Server on localhost:8080
	log.Print("Listening on localhost:8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
