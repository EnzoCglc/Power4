package controllers

import (
	"log"
	"net/http"
	"power4/models"
	"power4/utils"
	"strconv"
)

func GameDuo(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		reset(models.CurrentGame)
		models.CurrentGame.GameMode = "duo"
	}
	utils.Render(w, "gameBoard.html", models.CurrentGame)
}

func SwitchPlay(w http.ResponseWriter, r *http.Request) {
	exit := r.FormValue("exit")

	if exit == "reset" {
		reset(models.CurrentGame)
		utils.Render(w, "index.html", nil)
		return
	}

	colStr := r.FormValue("col")
	if colStr != "" {
		col, err := strconv.Atoi(colStr)
		if err != nil {
			log.Printf("invalid column %q: %v", colStr, err)
		} else {
			err = play(models.CurrentGame, col)
			if err != nil {
				log.Printf("invalid column %q: %v", colStr, err)
			}
		}
	}
	utils.Render(w, "gameBoard.html", models.CurrentGame)
}
