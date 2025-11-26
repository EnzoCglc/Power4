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
	level := parseBotLevel(r)
	initializeBotGame(level, r)

	username, err := getAuthenticatedUsername(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user, err := models.GetUserByUsername(username)
	if err != nil || user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data := buildGameBotData(user)
	utils.Render(w, "gameBoard.html", data)
}

// parseBotLevel extracts and validates bot difficulty level.
func parseBotLevel(r *http.Request) int {
	levelStr := r.URL.Query().Get("level")
	level, err := strconv.Atoi(levelStr)
	if err != nil || level < 1 || level > 5 {
		return 1
	}
	return level
}

// initializeBotGame resets and configures the game for bot mode.
func initializeBotGame(level int, r *http.Request) {
	reset(models.CurrentGame)
	models.CurrentGame.GameMode = "bot"
	models.CurrentGame.BotLvl = level

	rankedValue := getRankedValue(r)
	models.CurrentGame.Ranked = (rankedValue == "true")
}

// getRankedValue extracts ranked parameter from request.
func getRankedValue(r *http.Request) string {
	if r.Method == "POST" {
		return r.FormValue("ranked")
	}
	return r.URL.Query().Get("ranked")
}

// buildGameBotData constructs the game board data for rendering.
func buildGameBotData(user *models.User) map[string]any {
	winRate, totalGames := calculateUserStats(user)

	data := map[string]any{
		"CurrentGame": models.CurrentGame,
		"Player1":     user.Username,
		"User":        user,
		"WinRate":     winRate,
		"TotalGames":  totalGames,
	}

	if models.CurrentGame.Ranked {
		addBotStatsToData(data)
	} else {
		data["Player2"] = "Player 2"
	}

	return data
}

// calculateUserStats computes win rate and total games.
func calculateUserStats(user *models.User) (float64, int) {
	totalGames := user.Win + user.Losses
	var winRate float64
	if totalGames > 0 {
		winRate = (float64(user.Win) / float64(totalGames)) * 100
	}
	return winRate, totalGames
}

// addBotStatsToData adds bot player stats to game data.
func addBotStatsToData(data map[string]any) {
	player2, err := models.GetUserByUsername("player2")
	if err != nil || player2 == nil {
		log.Println("Player 2 not found in database, using default values")
		return
	}

	winRate2, totalGames2 := calculateUserStats(player2)
	data["Player2"] = player2.Username
	data["User2"] = player2
	data["WinRate2"] = winRate2
	data["TotalGames2"] = totalGames2
}

// SwitchPlayBot handles move submissions for bot mode games via AJAX.
func SwitchPlayBot(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Col int `json:"col"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		JSONError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	game := models.CurrentGame

	if err := play(game, request.Col); err != nil {
		JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if game.GameOver {
		sendGameResponse(w, game)
		return
	}

	executeBotMove(game)
	sendGameResponse(w, game)
}

// executeBotMove performs the bot's turn.
func executeBotMove(game *models.GridPage) {
	botCol := bot.BotMove(game, game.BotLvl, models.P2)
	if botCol != -1 {
		play(game, botCol)
	}
}

// sendGameResponse sends the game state as JSON response.
func sendGameResponse(w http.ResponseWriter, game *models.GridPage) {
	JSONSuccess(w, map[string]any{
		"game": game,
	})
}