package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"power4/models"
	"power4/utils"
)

// GameDuo initializes and renders a two-player (duo) game on the same device.
func GameDuo(w http.ResponseWriter, r *http.Request) {
	// Reset the game board and configure for duo mode
	reset(models.CurrentGame)
	models.CurrentGame.GameMode = "duo"
	models.CurrentGame.Ranked = false

	// Verify user is logged in
	cookie, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Fetch user data for Player 1
	username := cookie.Value
	user, err := models.GetUserByUsername(username)

	if err != nil || user == nil {
		log.Println("User not found: ", username)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Prepare data for the game board template
	data := map[string]interface{}{
		"CurrentGame": models.CurrentGame,
		"Player1":     user.Username,
	}

	utils.Render(w, "gameBoard.html", data)
	log.Println("Duo mod active for : ", user.Username)
}

// SwitchPlay handles move submissions for duo mode games via AJAX.
func SwitchPlay(w http.ResponseWriter, r *http.Request) {
	// Define the expected request structure
	var request struct {
		Col  int    `json:"col"`   // Column index where player wants to drop piece
		Exit string `json:"reset"` // "reset" to exit the game
	}

	// Parse the JSON request body
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		JSONError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Handle game reset/exit request
	if request.Exit == "reset" {
		reset(models.CurrentGame)
		utils.Render(w, "index.html", nil)
		return
	}

	// Process the player's move
	if request.Col >= 0 {
		err = play(models.CurrentGame, request.Col)
		if err != nil {
			log.Printf("Invalid move in column %d: %v", request.Col, err)
			return
		}

		log.Printf("Player %d played in column %d", models.CurrentGame.CurrentTurn, request.Col)
	}

	// Return the updated game state to the client
	JSONSuccess(w, map[string]interface{}{
		"game": models.CurrentGame,
	})
}

// JSONError sends a JSON error response to the client.
func JSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := map[string]interface{}{
		"success": false,
		"error":   message,
	}

	json.NewEncoder(w).Encode(response)
}

// JSONSuccess sends a JSON success response to the client.
func JSONSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"success": true,
		"data":    data,
	}

	json.NewEncoder(w).Encode(response)
}