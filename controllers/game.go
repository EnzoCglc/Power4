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
	log.Println("Duo mod active")
}

func GameSolo(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		reset(models.CurrentGame)
		models.CurrentGame.GameMode = "solo"

		lvl := r.FormValue("lvl")
		if lvl != "" {
			level, _ := strconv.Atoi(lvl)
			models.CurrentGame.BotLvl = level
		}
	}
	utils.Render(w, "gameBoard.html", models.CurrentGame)
	log.Println("Solo mod active")
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

			// Si c'est le mode solo et que le jeu n'est pas terminé, le bot joue immédiatement
			if models.CurrentGame.GameMode == "solo" && models.CurrentGame.CurrenctTurn == models.P2 && !models.CurrentGame.GameOver {
				playBot(models.CurrentGame)
				log.Println("Bot played")
			}
		}
	}
	utils.Render(w, "gameBoard.html", models.CurrentGame)
}
