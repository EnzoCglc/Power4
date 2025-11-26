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
	// Verify user is logged in by checking for username cookie
	cookie, err := r.Cookie("username")
	if err != nil {
		log.Println("User is not login")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Fetch user data from database
	username := cookie.Value
	user, err := models.GetUserByUsername(username)
	if err != nil {
		log.Println("Error fetching user from database:", err)
		http.Error(w, "Error loading data", http.StatusInternalServerError)
		return
	}

	// Handle case where user exists in cookie but not in database
	if user == nil {
		log.Println("User not exits in database:", username)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Calculate user statistics
	winRate := calculateWinRate(user.Win, user.Losses)
	totalGames := user.Win + user.Losses

	// Retrieve match history for this user
	history, err := models.GetHistoryByPlayer(username)
	if err != nil {
		log.Printf("Error to get history for %s : %s", username, err)
		return
	}

	// Prepare data for template rendering
	data := ProfilData{
		User:       user,
		WinRate:    winRate,
		TotalGames: totalGames,
		History:    history,
	}

	// Extract feedback messages from query parameters (used after redirects)
	if errorMsg := r.URL.Query().Get("error"); errorMsg != "" {
		data.ErrorMessage = errorMsg
	}
	if successMsg := r.URL.Query().Get("success"); successMsg != "" {
		data.SuccessMessage = successMsg
	}

	// Render the profile page with all user data
	utils.Render(w, "profil.html", data)

	log.Println("info de l'user ", user, "- WinRate:", winRate, "% - Total games:", totalGames)
}
