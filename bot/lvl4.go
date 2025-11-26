package bot

import (
	"power4/models"
)

// Lvl4Bot implements minimax AI with alpha-beta pruning at depth 6.
func Lvl4Bot(game *models.GridPage, player int) int {
	depth := 6
	validMoves := GetValideMoves(game)
	opponent := GetNextPlayer(player)

	// Check for immediate win/block before deep search
	if immediateMove := checkImmediateMove(game, validMoves, player, opponent); immediateMove != -1 {
		return immediateMove
	}

	return calculateBestMove(game, validMoves, player, depth)
}
