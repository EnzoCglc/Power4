package controllers

import (
	"net/http"
	"power4/utils"
)

func Home(w http.ResponseWriter, r *http.Request) {
	username := ""
	cookie, err := r.Cookie("username")
	if err == nil {
		username = cookie.Value
	}
	utils.Render(w, "index.html", username)
}
