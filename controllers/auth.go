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
		utils.Render(w, "registerPage.html", "Passwords are not identical")
		return
	}

	exists, err := verifExists(username)

	if err != nil {
		log.Printf("Error checking if user exists: %v", err)
		utils.Render(w, "registerPage.html", "Error to load Database")
		return
	}

	if exists {
		log.Println("User already exists:", username)
		utils.Render(w, "registerPage.html", "Username already taken")
		return
	}

	err = createUser(username, password)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		utils.Render(w, "registerPage.html", "Error creating user")
		return
	}

	log.Println("Nouveau compte accepté :", username)
	utils.Render(w, "loginPage.html", "Account created successfully! Please log in.")
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
		utils.Render(w, "loginPage.html", "Invalid username or password")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "username",
		Value: username,
		Path:  "/",
	})

	log.Println("Login Success for:", username)
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

func NewPassword(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	old_password := r.FormValue("old_password")
	new_password1 := r.FormValue("new_password1")
	new_password2 := r.FormValue("new_password2")

	err := login(username, old_password)
	if err != nil {
		log.Printf("Current password incorrect for user %s: %v", username, err)
		http.Redirect(w, r, "/profil?error=current_password_incorrect", http.StatusSeeOther)
		return
	}

	if new_password1 != new_password2 {
		log.Printf("New passwords do not match for user %s", username)
		http.Redirect(w, r, "/profil?error=passwords_do_not_match", http.StatusSeeOther)
		return
	}

	if len(new_password1) < 8 {
		log.Printf("New password too short for user %s", username)
		http.Redirect(w, r, "/profil?error=password_too_short", http.StatusSeeOther)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(new_password1), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password for user %s: %v", username, err)
		http.Redirect(w, r, "/profil?error=internal_error", http.StatusSeeOther)
		return
	}

	err = models.UpdatePassword(username, string(hash))
	if err != nil {
		log.Printf("Error updating password in database for user %s: %v", username, err)
		http.Redirect(w, r, "/profil?error=database_error", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/profil?success=password_updated", http.StatusSeeOther)
}