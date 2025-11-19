package controllers

import (
	"log"
	"net/http"
	"power4/models"
	"power4/utils"

	"golang.org/x/crypto/bcrypt"
)

type Auth_Data struct {
	ErrorMessage   string
	SuccessMessage string
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	data := Auth_Data{}

	// Check for error or success messages in query params
	if errorMsg := r.URL.Query().Get("error"); errorMsg != "" {
		data.ErrorMessage = errorMsg
	}
	if successMsg := r.URL.Query().Get("success"); successMsg != "" {
		data.SuccessMessage = successMsg
	}
	utils.Render(w, "loginPage.html", data)
}

func RegisterPage(w http.ResponseWriter, r *http.Request) {
	data := Auth_Data{}

	// Check for error or success messages in query params
	if errorMsg := r.URL.Query().Get("error"); errorMsg != "" {
		data.ErrorMessage = errorMsg
	}
	if successMsg := r.URL.Query().Get("success"); successMsg != "" {
		data.SuccessMessage = successMsg
	}

	utils.Render(w, "registerPage.html", data)
}

func RegisterInfo(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	confirm := r.FormValue("confirm_password")

	if password != confirm {
		log.Println("Passwords are not identical")
		http.Redirect(w, r, "/signup?error=passwords_do_not_match", http.StatusSeeOther)
		return
	}

	if len(password) < 8 {
		log.Printf("Password too short for user %s", username)
		http.Redirect(w, r, "/signup?error=password_too_short", http.StatusSeeOther)
		return
	}

	exists, err := verifExists(username)

	if err != nil {
		log.Printf("Error checking if user exists: %v", err)
		http.Redirect(w, r, "/signup?error=database_error", http.StatusSeeOther)
		return
	}

	if exists {
		log.Println("User already exists:", username)
		http.Redirect(w, r, "/signup?error=username_already_exists", http.StatusSeeOther)
		return
	}

	err = createUser(username, password)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		http.Redirect(w, r, "/signup?error=internal_error", http.StatusSeeOther)
		return
	}

	log.Println("Nouveau compte accepté :", username)
	http.Redirect(w, r, "/signin?success=registration_successful", http.StatusSeeOther)
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
		http.Redirect(w, r, "/signin?error=invalid_login", http.StatusSeeOther)
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

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:  "username",
		Value: "",
		Path:  "/",
		MaxAge: -1,
	})
	http.Redirect(w,r,"/",http.StatusSeeOther)
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