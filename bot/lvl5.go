package bot

import (
	"power4/models"
)

// Lvl8Bot implements minimax with alpha-beta pruning (depth 8)
func Lvl5Bot(game *models.GridPage , player int) int {
	depth := 8
	validMoves := GetValideMoves(game)
	opponent := GetNextPlayer(player)

	// Check for immediate win/block before deep search
	if immediateMove := checkImmediateMove(game, validMoves, player, opponent); immediateMove != -1 {
		return immediateMove
	}

	return calculateBestMove(game, validMoves, player, depth)
}
