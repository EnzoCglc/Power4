package bot

import (
	"math/rand"
	"power4/models"
)

// Lvl1Bot implements the easiest AI difficulty using basic defensive logic.
func Lvl1Bot(game *models.GridPage, player int) int {
	validMoves := GetValideMoves(game)
	if len(validMoves) == 0 {
		return -1 // No moves available (board full)
	}

	opponent := GetNextPlayer(player)
	safeMoves := []int{}

	// Filter out moves that would allow opponent to win immediately
	for _, col := range validMoves {
		row := SimulateMove(game, col, player)

		isSafe := true
		// Check if opponent can win after this move
		for _, opCol := range GetValideMoves(game) {
			opRow := SimulateMove(game, opCol, opponent)

			if CheckWin(game, opponent, opCol, opRow) {
				isSafe = false // This move leads to opponent's win
			}
			UndoMove(game, opCol, opRow)
			if !isSafe {
				break // No need to check further
			}
		}
		UndoMove(game, col, row)

		if isSafe {
			safeMoves = append(safeMoves, col)
		}
	}

	// If all moves are unsafe, pick randomly from all moves
	if len(safeMoves) == 0 {
		return validMoves[rand.Intn(len(validMoves))]
	}

	// Pick randomly from safe moves
	return safeMoves[rand.Intn(len(safeMoves))]
}