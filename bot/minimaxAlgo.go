package bot

import (
	"math"
	"power4/models"
)

// checkImmediateMove returns winning move or blocking move, -1 otherwise
func checkImmediateMove(game *models.GridPage, validMoves []int, player int, opponent int) int {
	if winMove := findWinningMove(game, validMoves, player); winMove != -1 {
		return winMove
	}
	return findWinningMove(game, validMoves, opponent)
}

// findWinningMove tests each column for winning move, returns column or -1
func findWinningMove(game *models.GridPage, validMoves []int, player int) int {
	for _, col := range validMoves {
		row := SimulateMove(game, col, player)
		if CheckWin(game, player, col, row) {
			UndoMove(game, col, row)
			return col
		}
		UndoMove(game, col, row)
	}
	return -1
}

// calculateBestMove uses minimax with alpha-beta pruning to find optimal move
func calculateBestMove(game *models.GridPage, validMoves []int, player int, depth int) int {
	bestScore := math.MinInt32
	bestMove := -1
	alpha := math.MinInt32
	beta := math.MaxInt32

	for _, col := range validMoves {
		row := SimulateMove(game , col, player)
		score := minimax(game, depth-1, false, player, alpha, beta)
		UndoMove(game, col, row)

		if score > bestScore {
			bestScore = score
			bestMove = col
		}
		alpha = max(alpha, score)
	}

	return bestMove
}

// minimax recursively evaluates game states, alternating between max/min players
func minimax(game *models.GridPage, depth int , isMaximizing bool, player int, alpha int, beta int) int {
	if score := checkTerminalState(game, depth, player); score != math.MaxInt32 {
		return score
	}

	validMoves := GetValideMoves(game)

	if isMaximizing {
		return maximizingPlayer(game, depth, player, alpha, beta, validMoves)
	}
	return minimizingPlayer(game, depth, player, alpha, beta, validMoves)
}

// checkTerminalState returns score if game ended, depth=0, or board full; else MaxInt32
func checkTerminalState(game *models.GridPage, depth int, player int) int {
	if winner := checkWinner(game); winner != models.Empty {
		if winner == player {
			return 1000
		}
		return -1000
	}
	if depth == 0 || IsBoardFull(game) {
		return evaluateBoard(game, player)
	}
	return math.MaxInt32
}

// maximizingPlayer finds best score for AI, uses alpha-beta pruning
func maximizingPlayer(game *models.GridPage, depth int, player int, alpha int, beta int, validMoves []int) int {
	maxEval := math.MinInt32
	for _, col := range validMoves {
		row := SimulateMove(game, col, player)
		eval := minimax(game, depth-1, false, player, alpha, beta)
		UndoMove(game, col, row)

		maxEval = max(maxEval, eval)
		alpha = max(alpha, eval)
		if beta <= alpha { // Prune branch
			break
		}
	}
	return maxEval
}

// minimizingPlayer finds worst score for AI (opponent's best), uses alpha-beta pruning
func minimizingPlayer(game *models.GridPage, depth int, player int, alpha int, beta int, validMoves []int) int {
	minEval := math.MaxInt32
	opponent := GetNextPlayer(player)
	for _, col := range validMoves {
		row := SimulateMove(game, col, opponent)
		eval := minimax(game, depth-1, true, player, alpha, beta)
		UndoMove(game, col, row)

		minEval = min(minEval, eval)
		beta = min(beta, eval)
		if beta <= alpha { // Prune branch
			break
		}
	}
	return minEval
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// evaluateBoard calculates heuristic score based on alignments and center control
func evaluateBoard(game *models.GridPage, player int) int {
	score := 0
	opponent := GetNextPlayer(player)

	score += evaluateAllDirections(game, player, opponent)
	score += evaluateCenterControl(game, player, opponent) // Center columns (2,3) are strategic

	return score
}

// evaluateCenterControl gives bonus for controlling center columns
func evaluateCenterControl(game *models.GridPage, player int, opponent int) int {
	score := 0
	centerCols := []int{2, 3}

	for _, col := range centerCols {
		for row := 0; row < models.Rows; row++ {
			if game.Columns[col][row] == player {
				score += 3
			} else if game.Columns[col][row] == opponent {
				score -= 3
			}
		}
	}

	return score
}

// evaluateAllDirections scores all 4-piece windows in all directions
func evaluateAllDirections(game *models.GridPage, player int, opponent int) int {
	score := 0
	score += evaluateHorizontal(game, player, opponent)
	score += evaluateVertical(game, player, opponent)
	score += evaluateDiagonalAscending(game, player, opponent)
	score += evaluateDiagonalDescending(game, player, opponent)
	return score
}

// evaluateHorizontal scores all horizontal 4-piece windows
func evaluateHorizontal(game *models.GridPage, player int, opponent int) int {
	score := 0
	for row := 0; row < models.Rows; row++ {
		for col := 0; col < models.Cols-3; col++ {
			window := []int{
				game.Columns[col][row],
				game.Columns[col+1][row],
				game.Columns[col+2][row],
				game.Columns[col+3][row],
			}
			score += evaluateWindow(window, player, opponent)
		}
	}
	return score
}

// evaluateVertical scores all vertical 4-piece windows
func evaluateVertical(game *models.GridPage, player int, opponent int) int {
	score := 0
	for col := 0; col < models.Cols; col++ {
		for row := 0; row < models.Rows-3; row++ {
			window := []int{
				game.Columns[col][row],
				game.Columns[col][row+1],
				game.Columns[col][row+2],
				game.Columns[col][row+3],
			}
			score += evaluateWindow(window, player, opponent)
		}
	}
	return score
}

// evaluateDiagonalAscending scores all ascending diagonal 4-piece windows (/)
func evaluateDiagonalAscending(game *models.GridPage, player int, opponent int) int {
	score := 0
	for col := 0; col < models.Cols-3; col++ {
		for row := 0; row < models.Rows-3; row++ {
			window := []int{
				game.Columns[col][row],
				game.Columns[col+1][row+1],
				game.Columns[col+2][row+2],
				game.Columns[col+3][row+3],
			}
			score += evaluateWindow(window, player, opponent)
		}
	}
	return score
}

// evaluateDiagonalDescending scores all descending diagonal 4-piece windows (\)
func evaluateDiagonalDescending(game *models.GridPage, player int, opponent int) int {
	score := 0
	for col := 0; col < models.Cols-3; col++ {
		for row := 3; row < models.Rows; row++ {
			window := []int{
				game.Columns[col][row],
				game.Columns[col+1][row-1],
				game.Columns[col+2][row-2],
				game.Columns[col+3][row-3],
			}
			score += evaluateWindow(window, player, opponent)
		}
	}
	return score
}

// evaluateWindow scores a 4-piece window based on piece counts
func evaluateWindow(window []int, player int, opponent int) int {
	playerCount, opponentCount, emptyCount := countPieces(window, player, opponent)

	if playerCount > 0 && opponentCount > 0 { // Mixed window is useless
		return 0
	}

	return calculateWindowScore(playerCount, opponentCount, emptyCount)
}

func countPieces(window []int, player int, opponent int) (int, int, int) {
	playerCount := 0
	opponentCount := 0
	emptyCount := 0

	for _, cell := range window {
		if cell == player {
			playerCount++
		} else if cell == opponent {
			opponentCount++
		} else {
			emptyCount++
		}
	}

	return playerCount, opponentCount, emptyCount
}

// calculateWindowScore assigns strategic values to piece patterns
func calculateWindowScore(playerCount int, opponentCount int, emptyCount int) int {
	score := 0

	// Offensive scoring
	if playerCount == 4 {
		score += 100
	} else if playerCount == 3 && emptyCount == 1 {
		score += 10
	} else if playerCount == 2 && emptyCount == 2 {
		score += 5
	}

	// Defensive scoring (blocking threats)
	if opponentCount == 3 && emptyCount == 1 {
		score -= 80
	} else if opponentCount == 2 && emptyCount == 2 {
		score -= 4
	}

	return score
}

// checkWinner scans board for any winning player, returns player ID or Empty
func checkWinner(game *models.GridPage) int {
	for col := 0; col < models.Cols; col++ {
		for row := 0; row < models.Rows; row++ {
			if game.Columns[col][row] != models.Empty {
				player := game.Columns[col][row]
				if CheckWin(game, player, col, row) {
					return player
				}
			}
		}
	}
	return models.Empty
}