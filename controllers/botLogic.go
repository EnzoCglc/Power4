package controllers

import (
	"math/rand/v2"
	"power4/models"
)

func playBot(game *models.GridPage) {
	var col int
	switch game.BotLvl {
	case 1:
		col = randomMove(game)
	}
	play(game, col)
}

func randomMove(game *models.GridPage) int {
	colPlay:= []int{}

	for col := 0; col < models.Cols; col++ {
		if game.Columns[col][0] == models.Empty {
			colPlay = append(colPlay, col)
		}
	}
	random := rand.IntN(len(colPlay))
	return colPlay[random]
}
