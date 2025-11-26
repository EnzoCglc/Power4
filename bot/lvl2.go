package bot

import (
	"power4/models"
)

// Lvl2Bot implements minimax AI with alpha-beta pruning at depth 2.
func Lvl2Bot(game *models.GridPage, player int) int {
	depth := 2
	validMoves := GetValideMoves(game)
	opponent := GetNextPlayer(player)

	// Optimization: Check for immediate win or necessary block first
	// This saves computation time by avoiding minimax for obvious moves
	if immediateMove := checkImmediateMove(game, validMoves, player, opponent); immediateMove != -1 {
		return immediateMove
	}

	// Use minimax to find the best move by looking 2 turns ahead
	return calculateBestMove(game, validMoves, player, depth)
}
