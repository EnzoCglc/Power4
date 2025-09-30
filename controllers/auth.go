package controllers

import (
	"log"
	"net/http"
	"power4/models"
	"power4/utils"
)

func LoginPage(w http.ResponseWriter, r *http.Request) {
	utils.Render(w, "loginPage.html", nil)
}

func RegisterPage(w http.ResponseWriter, r *http.Request) {
	utils.Render(w, "registerPage.html", nil)
}

func RegisterInfo(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	confirm := r.FormValue("confirm_password")

	if password != confirm {
		log.Println("MDP INCORRECT")
		utils.Render(w, "registerPage.html", "Password don't match")
		return
	}
	verifDataBase(w, username)
	log.Println("Nouveau compte accepté :", username)
	utils.Render(w, "registerPage.html", nil)
}

func verifDataBase(w http.ResponseWriter, username string) {
	db, err := models.LoadDB("database/db.json")
	if err != nil {
		http.Error(w, "Error to load Database", http.StatusInternalServerError)
		return
	}
	for _, u := range db.User {
		if u.Username == username {
			log.Println("Username already use : ", username)
			utils.Render(w, "registerPage.html", "Username already use ")
			return
		}
	}
}
