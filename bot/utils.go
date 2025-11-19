package bot

import (
	"power4/models"
)

func GetValideMoves(game *models.GridPage) []int {
	validMoves := []int{}
	for col := 0; col < models.Cols; col++ {
		if game.Columns[col][0] == models.Empty {
			validMoves = append(validMoves, col)
		}
	}
	return validMoves
}

func findAvailableRow(cols [][]int, col int) int {
	for row := models.Rows - 1; row >= 0; row-- {
		if cols[col][row] == models.Empty {
			return row
		}
	}
	return -1
}

func SimulateMove(game *models.GridPage, col int, player int) int {
	row := findAvailableRow(game.Columns, col)
	if row != -1 {
		game.Columns[col][row] = player
	}
	return row
}

func UndoMove(game *models.GridPage, col int , row int){
	if row != -1 {
		game.Columns[col][row] = models.Empty
	}
}

func verifWin(cols [][]int, player int, col int, row int) bool {
	grid := [][2]int{
		{1, 0},  // horizontal
		{0, 1},  // vertical
		{1, 1},  // diagonal \
		{1, -1}, // diagonal /
	}

	for _, g := range grid {
		count := 1
		count += countDirection(cols, player, col, row, g[0], g[1])
		count += countDirection(cols, player, col, row, -g[0], -g[1])

		if count >= 4 {
			return true
		}
	}
	return false
}

func countDirection(cols [][]int, player int, col int, row int, dc int, dr int) int {
	c := col + dc
	r := row + dr
	count := 0
	for c >= 0 && c < models.Cols && r >= 0 && r < models.Rows && cols[c][r] == player {
		count += 1
		c += dc
		r += dr
	}
	return count
}

func CheckWin(game *models.GridPage, player int, col int, row int) bool {
	if row == -1 {
		return false
	}
	return verifWin(game.Columns, player, col, row)
}

func GetNextPlayer(currentplayer int) int {
	if currentplayer == models.P1 {
		return models.P2
	}
	return models.P1
}

func IsBoardFull(game *models.GridPage) bool {
	for col := 0; col < models.Cols; col++ {
		if game.Columns[col][0] == models.Empty {
			return false
		}
	}
	return true
}