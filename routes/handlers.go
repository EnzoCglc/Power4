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
	http.HandleFunc("/game/duo", controllers.GameDuo)
	http.HandleFunc("/game", controllers.SwitchPlay)

	//Result game
	http.HandleFunc("/game/result", controllers.GameResult)

	//Connect LoginPage
	http.HandleFunc("/signin", controllers.LoginPage)
	http.HandleFunc("/login", controllers.LoginInfo)

	//Connect RegisterPage
	http.HandleFunc("/signup", controllers.RegisterPage)
	http.HandleFunc("/register", controllers.RegisterInfo)

	//Connect Profil Page
	http.HandleFunc("/profil", controllers.Profil)
	http.HandleFunc("/profil/update-password", controllers.NewPassword)

	//Logout
	http.HandleFunc("/logout", controllers.Logout)

	// Use FileServer to serve static assets like .png or css
	fs := http.FileServer(http.Dir("./assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
}
