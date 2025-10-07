package routes

import (
	"net/http"
	"power4/controllers"
)

func SetupRoutes() {
	// Connect landing Page
	http.HandleFunc("/", controllers.Home)

	//Connect Gamemode

	http.HandleFunc("/gamemode", controllers.GameMode)

	//Connect game function
	http.HandleFunc("/game", controllers.SwitchPlay)

	// Use FileServer to serve static assets like .png or css
	fs := http.FileServer(http.Dir("./assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
}
