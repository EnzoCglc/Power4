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

func RegisterPage(w http.ResponseWriter, r *http.Request) {
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
		log.Printf("Error checking if user exists: %v", err)
		http.Error(w, "Error to load Database", http.StatusInternalServerError)
		return
	}

	if exists {
		log.Println("User already use")
		utils.Render(w, "registerPage.html", nil)
		return
	}

	err = createUser(username, password)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	log.Println("Nouveau compte accepté :", username)
	utils.Render(w, "loginPage.html", nil)
}

func verifExists(username string) (bool, error) {
	return models.UserExists(username)
}

func createUser(username, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return models.CreateUser(username, string(hash))
}

func LoginInfo(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	err := login(username, password)

	if err != nil {
		log.Printf("Login error for user %s: %v", username, err)
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name: "username",
		Value: username,
		Path: "/",
	})

	log.Println("Login Succes for :", username)
	utils.Render(w, "index.html", username)
}

func login(username, password string) error {
	user, err := models.GetUserByUsername(username)
	if err != nil {
		return err
	}
	if user == nil {
		return bcrypt.ErrMismatchedHashAndPassword
	}

	return bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
}
