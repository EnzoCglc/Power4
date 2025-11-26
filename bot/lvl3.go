package bot

import (
	"power4/models"
)

// Lvl3Bot implements minimax AI with alpha-beta pruning at depth 4.
func Lvl3Bot(game *models.GridPage, player int) int {
	depth := 4
	validMoves := GetValideMoves(game)
	opponent := GetNextPlayer(player)

	// Check for immediate win/block before deep search
	if immediateMove := checkImmediateMove(game, validMoves, player, opponent); immediateMove != -1 {
		return immediateMove
	}

	return calculateBestMove(game, validMoves, player, depth)
}
