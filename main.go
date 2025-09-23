package main

import (
	"log"
	"net/http"
	"power4/controllers"
)

func main() {
	// Connect landing Page
	http.HandleFunc("/", controllers.Home)

	// Use FileServer to serve static assets like .png or css
	fs := http.FileServer(http.Dir("./assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	// Start Server on localhost:8080
	log.Print("Listening on localhost:8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
