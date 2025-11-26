package controllers

import (
	"log"
	"net/http"
	"power4/models"
	"power4/utils"
)

// ProfilData contains all the information needed to render a user's profile page.
// This struct is passed to the template engine to display user statistics,
// match history, and any messages (errors or success notifications).
type ProfilData struct {
	User           *models.User      // User account information (username, ELO, wins, losses)
	WinRate        float64           // Calculated win percentage (0-100)
	TotalGames     int               // Total number of games played (wins + losses)
	History        []models.History  // List of past matches involving this user
	ErrorMessage   string            // Error message to display (e.g., "password_too_short")
	SuccessMessage string            // Success message to display (e.g., "password_updated")
}

// calculateWinRate computes the win percentage based on wins and losses.
//
// Parameters:
//   - wins: Number of games won
//   - losses: Number of games lost
//
// Returns:
//   - Win rate as a percentage (0.0 to 100.0)
//   - Returns 0.0 if no games have been played
//
// This function handles the edge case where a new user has played no games,
// preventing division by zero errors.
func calculateWinRate(wins, losses int) float64 {
	totalGames := wins + losses
	if totalGames == 0 {
		return 0.0
	}
	return (float64(wins) / float64(totalGames)) * 100
}

// Profil handles the user profile page request, displaying stats and match history.
func Profil(w http.ResponseWriter, r *http.Request) {
	username, err := getAuthenticatedUsername(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user, err := fetchUserProfile(username)
	if err != nil {
		handleProfilError(w, r, err)
		return
	}

	history, err := models.GetHistoryByPlayer(username)
	if err != nil {
		log.Printf("Error to get history for %s : %s", username, err)
		return
	}

	data := buildProfilData(user, history, r)
	utils.Render(w, "profil.html", data)

	log.Println("info de l'user ", user, "- WinRate:", data.WinRate, "% - Total games:", data.TotalGames)
}

// getAuthenticatedUsername retrieves the logged-in username from cookie.
func getAuthenticatedUsername(r *http.Request) (string, error) {
	cookie, err := r.Cookie("username")
	if err != nil {
		log.Println("User is not login")
		return "", err
	}
	return cookie.Value, nil
}

// fetchUserProfile retrieves user from database.
func fetchUserProfile(username string) (*models.User, error) {
	user, err := models.GetUserByUsername(username)
	if err != nil {
		log.Println("Error fetching user from database:", err)
		return nil, err
	}

	if user == nil {
		log.Println("User not exits in database:", username)
		return nil, http.ErrNoCookie
	}

	return user, nil
}

// handleProfilError handles errors during profile loading.
func handleProfilError(w http.ResponseWriter, r *http.Request, err error) {
	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		http.Error(w, "Error loading data", http.StatusInternalServerError)
	}
}

// buildProfilData constructs the profile page data structure.
func buildProfilData(user *models.User, history []models.History, r *http.Request) ProfilData {
	winRate := calculateWinRate(user.Win, user.Losses)
	totalGames := user.Win + user.Losses

	data := ProfilData{
		User:       user,
		WinRate:    winRate,
		TotalGames: totalGames,
		History:    history,
	}

	if errorMsg := r.URL.Query().Get("error"); errorMsg != "" {
		data.ErrorMessage = errorMsg
	}
	if successMsg := r.URL.Query().Get("success"); successMsg != "" {
		data.SuccessMessage = successMsg
	}

	return data
}
