package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"power4/bot"
	"power4/models"
	"power4/utils"
	"strconv"
)

// GameBot initializes and renders a game against an AI opponent.
func GameBot(w http.ResponseWriter, r *http.Request) {
	// Parse and validate bot difficulty level
	levelStr := r.URL.Query().Get("level")
	level, err := strconv.Atoi(levelStr)
	if err != nil || level < 1 || level > 5 {
		level = 1 // Default to easiest level if invalid
	}

	// Reset and configure game for bot mode
	reset(models.CurrentGame)
	models.CurrentGame.GameMode = "bot"
	models.CurrentGame.BotLvl = level

	// Determine ranked status (supports both POST and GET)
	var rankedValue string

	switch r.Method {
	case "POST":
		rankedValue = r.FormValue("ranked")
	case "GET":
		rankedValue = r.URL.Query().Get("ranked")
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Set ranked flag (only "true" string sets it to true)
	if rankedValue == "true" {
		models.CurrentGame.Ranked = true
	} else {
		models.CurrentGame.Ranked = false
	}

	// Verify user authentication
	cookie, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Fetch Player 1 (human player) data
	username := cookie.Value
	user, err := models.GetUserByUsername(username)
	if err != nil || user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Calculate Player 1 statistics
	totalGames := user.Win + user.Losses
	var winRate float64
	if totalGames > 0 {
		winRate = (float64(user.Win) / float64(totalGames)) * 100
	}

	// Prepare base game data
	data := map[string]interface{}{
		"CurrentGame": models.CurrentGame,
		"Player1":     user.Username,
		"User":        user,
		"WinRate":     winRate,
		"TotalGames":  totalGames,
	}

	// For ranked games, fetch bot's persistent user stats (stored as "player2")
	if models.CurrentGame.Ranked {
		player2, err := models.GetUserByUsername("player2")
		if err != nil || player2 == nil {
			log.Println("Player 2 not found in database, using default values")
		} else {
			// Calculate bot's statistics
			totalGames2 := player2.Win + player2.Losses
			var winRate2 float64
			if totalGames2 > 0 {
				winRate2 = (float64(player2.Win) / float64(totalGames2)) * 100
			}
			data["Player2"] = player2.Username
			data["User2"] = player2
			data["WinRate2"] = winRate2
			data["TotalGames2"] = totalGames2
		}
	} else {
		// For unranked games, use generic name
		data["Player2"] = "Player 2"
	}

	utils.Render(w, "gameBoard.html", data)
}

// SwitchPlayBot handles move submissions for bot mode games via AJAX.
func SwitchPlayBot(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Col int `json:"col"` // Column where player wants to drop piece
	}

	// Parse player's move from JSON
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		JSONError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	game := models.CurrentGame

	// Step 1: Execute player's move
	err = play(game, request.Col)
	if err != nil {
		JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if player won or board is full
	if game.GameOver {
		JSONSuccess(w, map[string]interface{}{
			"game": game,
		})
		return
	}

	// Step 2: Execute bot's move (P2)
	botCol := bot.BotMove(game, game.BotLvl, models.P2)
	if botCol != -1 {
		play(game, botCol)
	}

	// Return updated game state after both moves
	JSONSuccess(w, map[string]interface{}{
		"game": game,
	})
}