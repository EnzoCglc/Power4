package controllers

import (
	"encoding/json"
	"net/http"
	"power4/models"
	"power4/utils"
	"power4/bot"
	"strconv"
	"log"
)

func GameBot(w http.ResponseWriter, r *http.Request) {
	levelStr := r.URL.Query().Get("level")
	level, err := strconv.Atoi(levelStr)
	if err != nil || level < 1 || level > 5 {
		level = 1
	}

	reset(models.CurrentGame)
	models.CurrentGame.GameMode = "bot"
	models.CurrentGame.BotLvl = level

	var rankedValue string

	switch r.Method {
	case "POST":
		rankedValue = r.FormValue("ranked")
	case "GET":
		rankedValue = r.URL.Query().Get("ranked")
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	if rankedValue == "true" {
		models.CurrentGame.Ranked = true
	} else {
		models.CurrentGame.Ranked = false
	}

	cookie, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	username := cookie.Value
	user, err := models.GetUserByUsername(username)
	if err != nil || user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Calculer le winrate du Player 1
	totalGames := user.Win + user.Losses
	var winRate float64
	if totalGames > 0 {
		winRate = (float64(user.Win) / float64(totalGames)) * 100
	}

	data := map[string]interface{}{
		"CurrentGame": models.CurrentGame,
		"Player1":     user.Username,
		"User":        user,
		"WinRate":     winRate,
		"TotalGames":  totalGames,
	}

	if models.CurrentGame.Ranked {
		player2, err := models.GetUserByUsername("player2")
		if err != nil || player2 == nil {
			log.Println("Player 2 not found in database, using default values")
		} else {
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
		data["Player2"] = "Player 2"
	}
	utils.Render(w, "gameBoard.html", data)
}

func SwitchPlayBot(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Col int `json:"col"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		JSONError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	game := models.CurrentGame

	// 1. Coup du joueur
	err = play(game, request.Col)
	if err != nil {
		JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Vérifier si le joueur a gagné
	if game.GameOver {
		JSONSuccess(w, map[string]interface{}{
			"game": game,
		})
		return
	}

	// 2. Coup du bot
	botCol := bot.BotMove(game, game.BotLvl, models.P2)
	if botCol != -1 {
		play(game, botCol)
	}

	JSONSuccess(w, map[string]interface{}{
		"game": game,
	})
}