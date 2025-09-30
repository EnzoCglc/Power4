package controllers

import (
	"log"
	"net/http"
	"power4/models"
	"power4/utils"

	"golang.org/x/crypto/bcrypt"
)

func LoginPage(w http.ResponseWriter, r *http.Request) {
	utils.Render(w, "loginPage.html", nil)
}

func RegisauterPage(w http.ResponseWriter, r *http.Request) {
	utils.Render(w, "registerPage.html", nil)
}

func RegisterInfo(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	confirm := r.FormValue("confirm_password")

	if password != confirm {
		utils.Render(w, "registerPage.html", "Passwords don't match")
		return
	}

	exists, err := verifExists(username)

	if err != nil {
		http.Error(w, "Error to load Database", http.StatusInternalServerError)
		return
	}

	if exists {
		log.Println("User already use")
		utils.Render(w, "registerPage.html", nil)
		return
	}

	createUser(username, password)

	log.Println("Nouveau compte accepté :", username)
	utils.Render(w, "loginPage.html", nil)
}

func verifExists(username string) (bool, error) {
	db, err := models.LoadDB("database/db.json")

	if err != nil {
		return false, err
	}

	for _, u := range db.Users {
		if u.Username == username {
			return true, nil
		}
	}
	return false, nil
}

func createUser(username, password string) error {
	db, err := models.LoadDB("database/db.json")

	if err != nil {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	newUser := models.User{
		ID:           len(db.Users) + 1,
		Username:     username,
		PasswordHash: string(hash),
		Elo:          1000,
		Win:          0,
		Losses:       0,
	}
	db.Users = append(db.Users, newUser)
	return models.SaveDB("database/db.json", db)
}
