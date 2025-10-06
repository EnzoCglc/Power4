package controllers

import (
	"net/http"
	"power4/utils"
)

func GameMode(w http.ResponseWriter, r *http.Request) {
	utils.Render(w, "gamemode.html", nil)
}
