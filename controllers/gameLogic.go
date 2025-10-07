package controllers

import (
	"log"
	"power4/models"
)

func play(game *models.GridPage, col int) {
	player := game.CurrenctTurn

	for row := models.Rows - 1; row >= 0; row-- {
		if game.Columns[col][row] == models.Empty {
			game.Columns[col][row] = player
			if verifWin(game.Columns, player, col, row) {
				log.Printf("Player %d win", player)
				return
			} else {
				if player == models.P1 {
					game.CurrenctTurn = models.P2
				} else {
					game.CurrenctTurn = models.P1
				}
				return
			}
		}
	}
}

func reset(game *models.GridPage) {
	for i := 0; i < models.Cols; i++ {
		game.Columns[i] = make([]int, models.Rows)
	}
	game.CurrenctTurn = models.P1
}