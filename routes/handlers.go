package routes

import (
	"net/http"
	"power4/controllers"
)

func SetupRoutes() {
	// Main pages
	http.HandleFunc("/", controllers.Home)                  // Landing page
	http.HandleFunc("/gamemode", controllers.GameMode)      // Game mode selection

	// Duo mode (local multiplayer)
	http.HandleFunc("/game/duo", controllers.GameDuo)       // Initialize duo game
	http.HandleFunc("/game", controllers.SwitchPlay)        // Process moves in duo mode

	// Bot mode (single player vs AI)
	http.HandleFunc("/game/bot", controllers.GameBot)       // Initialize bot game
	http.HandleFunc("/game/bot/play", controllers.SwitchPlayBot)  // Process moves in bot mode

	// Game results API
	http.HandleFunc("/game/result", controllers.GameResult) // Process game completion and ELO

	// Authentication routes
	http.HandleFunc("/signin", controllers.LoginPage)       // Display login form
	http.HandleFunc("/login", controllers.LoginInfo)        // Process login submission
	http.HandleFunc("/signup", controllers.RegisterPage)    // Display registration form
	http.HandleFunc("/register", controllers.RegisterInfo)  // Process registration submission
	http.HandleFunc("/logout", controllers.Logout)          // Clear session and logout

	// User profile routes
	http.HandleFunc("/profil", controllers.Profil)          // Display user profile
	http.HandleFunc("/profil/update-password", controllers.NewPassword)  // Change password

	// Static file server for CSS, JS, images, etc.
	// Serves files from ./assets directory at /assets/ URL path
	fs := http.FileServer(http.Dir("./assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
}
