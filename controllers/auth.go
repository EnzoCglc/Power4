package controllers

import (
	"log"
	"net/http"
	"power4/models"
	"power4/utils"

	"golang.org/x/crypto/bcrypt"
)

// Auth_Data contains messages to display on authentication pages.
type Auth_Data struct {
	ErrorMessage   string // Error message key (e.g., "passwords_do_not_match")
	SuccessMessage string // Success message key (e.g., "registration_successful")
}

// LoginPage renders the login page with optional error/success messages.
func LoginPage(w http.ResponseWriter, r *http.Request) {
	data := Auth_Data{}

	// Extract feedback messages from URL query parameters
	if errorMsg := r.URL.Query().Get("error"); errorMsg != "" {
		data.ErrorMessage = errorMsg
	}
	if successMsg := r.URL.Query().Get("success"); successMsg != "" {
		data.SuccessMessage = successMsg
	}
	utils.Render(w, "loginPage.html", data)
}

// RegisterPage renders the registration page with optional error/success messages.
func RegisterPage(w http.ResponseWriter, r *http.Request) {
	data := Auth_Data{}

	// Extract feedback messages from URL query parameters
	if errorMsg := r.URL.Query().Get("error"); errorMsg != "" {
		data.ErrorMessage = errorMsg
	}
	if successMsg := r.URL.Query().Get("success"); successMsg != "" {
		data.SuccessMessage = successMsg
	}

	utils.Render(w, "registerPage.html", data)
}

// RegisterInfo processes the registration form submission and creates a new user account.
func RegisterInfo(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	confirm := r.FormValue("confirm_password")

	if err := validateRegistration(username, password, confirm); err != nil {
		http.Redirect(w, r, "/signup?error="+err.Error(), http.StatusSeeOther)
		return
	}

	if err := createUser(username, password); err != nil {
		log.Printf("Error creating user: %v", err)
		http.Redirect(w, r, "/signup?error=internal_error", http.StatusSeeOther)
		return
	}

	log.Println("Nouveau compte accepté :", username)
	http.Redirect(w, r, "/signin?success=registration_successful", http.StatusSeeOther)
}

// validateRegistration validates registration input and checks for existing users.
func validateRegistration(username, password, confirm string) error {
	if password != confirm {
		log.Println("Passwords are not identical")
		return &validationError{"passwords_do_not_match"}
	}

	if len(password) < 8 {
		log.Printf("Password too short for user %s", username)
		return &validationError{"password_too_short"}
	}

	exists, err := verifExists(username)
	if err != nil {
		log.Printf("Error checking if user exists: %v", err)
		return &validationError{"database_error"}
	}

	if exists {
		log.Println("User already exists:", username)
		return &validationError{"username_already_exists"}
	}

	return nil
}

type validationError struct {
	key string
}

func (e *validationError) Error() string {
	return e.key
}

// verifExists checks if a username already exists in the database.
func verifExists(username string) (bool, error) {
	return models.UserExists(username)
}

// createUser creates a new user account with a securely hashed password.
func createUser(username, password string) error {
	// Hash the password using bcrypt for secure storage
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Store the user with the hashed password
	return models.CreateUser(username, string(hash))
}

// LoginInfo processes the login form submission and authenticates the user.
func LoginInfo(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Verify credentials
	err := login(username, password)
	if err != nil {
		log.Printf("Login error for user %s: %v", username, err)
		http.Redirect(w, r, "/signin?error=invalid_login", http.StatusSeeOther)
		return
	}

	// Create session cookie
	http.SetCookie(w, &http.Cookie{
		Name:  "username",
		Value: username,
		Path:  "/",
	})

	log.Println("Login Success for:", username)
	utils.Render(w, "index.html", username)
}

// login verifies user credentials against the database.
func login(username, password string) error {
	// Fetch user from database
	user, err := models.GetUserByUsername(username)
	if err != nil {
		return err
	}

	// User not found
	if user == nil {
		return bcrypt.ErrMismatchedHashAndPassword
	}

	// Verify password using bcrypt (time-constant comparison)
	return bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
}

// Logout clears the user's session and redirects to the home page.
func Logout(w http.ResponseWriter, r *http.Request) {
	// Delete the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "username",
		Value:  "",
		Path:   "/",
		MaxAge: -1, // MaxAge < 0 means delete cookie immediately
	})

	// Redirect to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// NewPassword handles password change requests from the user profile page.
func NewPassword(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	oldPassword := r.FormValue("old_password")
	newPassword1 := r.FormValue("new_password1")
	newPassword2 := r.FormValue("new_password2")

	if err := validatePasswordChange(username, oldPassword, newPassword1, newPassword2); err != nil {
		http.Redirect(w, r, "/profil?error="+err.Error(), http.StatusSeeOther)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword1), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password for user %s: %v", username, err)
		http.Redirect(w, r, "/profil?error=internal_error", http.StatusSeeOther)
		return
	}

	if err := models.UpdatePassword(username, string(hash)); err != nil {
		log.Printf("Error updating password in database for user %s: %v", username, err)
		http.Redirect(w, r, "/profil?error=database_error", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/profil?success=password_updated", http.StatusSeeOther)
}

// validatePasswordChange validates password change request.
func validatePasswordChange(username, oldPassword, newPassword1, newPassword2 string) error {
	if err := login(username, oldPassword); err != nil {
		log.Printf("Current password incorrect for user %s: %v", username, err)
		return &validationError{"current_password_incorrect"}
	}

	if newPassword1 != newPassword2 {
		log.Printf("New passwords do not match for user %s", username)
		return &validationError{"passwords_do_not_match"}
	}

	if len(newPassword1) < 8 {
		log.Printf("New password too short for user %s", username)
		return &validationError{"password_too_short"}
	}

	return nil
}