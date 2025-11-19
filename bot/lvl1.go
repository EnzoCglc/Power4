package bot

import (
	"power4/models"
	"math/rand"
)

func Lvl1Bot(game *models.GridPage , player int) int {
	validMoves := GetValideMoves(game)
	if len(validMoves) == 0 {
		return -1 
	}

	opponent := GetNextPlayer(player)
	safeMoves := []int{}

	for _, col := range validMoves {
		row := SimulateMove(game, col, player)

		isSafe := true
		for _, opCol := range GetValideMoves(game) {
			opRow := SimulateMove(game , opCol, opponent)

			if CheckWin(game, opponent, opCol ,opRow) {
				isSafe = false
			}
			UndoMove(game, opCol, opRow)
			if !isSafe {
				break
			}
		}
		UndoMove(game, col, row)
		if isSafe {
			safeMoves = append(safeMoves, col)
		}
	}

	if len(safeMoves) == 0 {
		return validMoves[rand.Intn(len(validMoves))]
	}
	return safeMoves[rand.Intn(len(safeMoves))]
}