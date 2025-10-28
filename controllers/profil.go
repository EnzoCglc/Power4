package controllers

import (
	"log"
	"net/http"
	"power4/models"
	"power4/utils"
)

func Profil(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("username")
	if err != nil {
		log.Println("User is not login")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	username := cookie.Value
	user, err := models.GetUserByUsername(username)
	if err != nil {
		log.Println("Error fetching user from database:", err)
		http.Error(w, "Error loading data", http.StatusInternalServerError)
		return
	}

	if user == nil {
		log.Println("User not exits in database:", username)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	utils.Render(w, "profil.html", user)

	log.Println("info de l'user ", user)
}
