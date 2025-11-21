package controllers

import (
	"encoding/json"
	"net/http"
	"power4/models"
	"power4/utils"
	"power4/bot"
	"strconv"
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
	models.CurrentGame.Ranked = false

	cookie, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	username := cookie.Value
	user, err := models.GetUserByUsername(username)
	if err != nil || user == nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	data := map[string]interface{}{
		"CurrentGame": models.CurrentGame,
		"Player1":     user.Username,
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