package controllers

import (
	"net/http"
	"power4/utils"
)

// GameMode handles the game mode selection page request.
func GameMode(w http.ResponseWriter, r *http.Request) {
	utils.Render(w, "gamemode.html", nil)
}
