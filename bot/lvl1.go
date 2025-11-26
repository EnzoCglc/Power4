package bot

import (
	"math/rand"
	"power4/models"
)

// Lvl1Bot implements the easiest AI difficulty using basic defensive logic.
func Lvl1Bot(game *models.GridPage, player int) int {
	validMoves := GetValideMoves(game)
	if len(validMoves) == 0 {
		return -1
	}

	opponent := GetNextPlayer(player)
	safeMoves := filterSafeMoves(game, validMoves, player, opponent)

	if len(safeMoves) == 0 {
		return validMoves[rand.Intn(len(validMoves))]
	}

	return safeMoves[rand.Intn(len(safeMoves))]
}

// filterSafeMoves filters out moves that would allow opponent to win.
func filterSafeMoves(game *models.GridPage, validMoves []int, player, opponent int) []int {
	safeMoves := []int{}

	for _, col := range validMoves {
		if isMoveSafe(game, col, player, opponent) {
			safeMoves = append(safeMoves, col)
		}
	}

	return safeMoves
}

// isMoveSafe checks if a move doesn't lead to opponent's immediate win.
func isMoveSafe(game *models.GridPage, col, player, opponent int) bool {
	row := SimulateMove(game, col, player)
	defer UndoMove(game, col, row)

	for _, opCol := range GetValideMoves(game) {
		opRow := SimulateMove(game, opCol, opponent)
		isWin := CheckWin(game, opponent, opCol, opRow)
		UndoMove(game, opCol, opRow)

		if isWin {
			return false
		}
	}

	return true
}