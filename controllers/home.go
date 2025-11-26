package controllers

import (
	"net/http"
	"power4/utils"
)

// Home handles the main landing page of the application.
func Home(w http.ResponseWriter, r *http.Request) {
	username := ""

	// Attempt to retrieve the username from the session cookie
	cookie, err := r.Cookie("username")
	if err == nil {
		username = cookie.Value
	}

	// Render the home page with the username (empty if not logged in)
	utils.Render(w, "index.html", username)
}
