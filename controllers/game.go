package controllers

import (
	"log"
	"net/http"
	m "power4/models"
	"strconv"
)

func SwitchPlay(w http.ResponseWriter, r *http.Request) {
	action := r.FormValue("play")
	log.Println(action)

	colStr := r.FormValue("col")
	if colStr != "" {
		col, err := strconv.Atoi(colStr)
		if err != nil {
			log.Printf("invalid column %q: %v", colStr, err)
		} else {
			play(m.CurrentGame, col)
		}
	}
	render(w, "gameBoard.html", m.CurrentGame)
}

func play(game *m.GridPage, col int) {
	player := game.CurrenctTurn

	for row := m.Rows - 1; row >= 0; row-- {
		if game.Columns[col][row] == m.Empty {
			game.Columns[col][row] = player
			log.Println("Value de la grille " , game.Columns)
			if player == m.P1 {
				game.CurrenctTurn = m.P2
			} else {
				game.CurrenctTurn = m.P1
			}
			return
		}
	}
}
