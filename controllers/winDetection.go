package controllers

import (
	"power4/models"
)

// VerifWin checks if a player has won the game by getting 4 pieces in a row.
func VerifWin(cols [][]int, player int, col int, row int) bool {
	// Define direction vectors: [deltaCol, deltaRow]
	grid := [][2]int{
		{1, 0},  // horizontal: check left and right
		{0, 1},  // vertical: check up and down
		{1, 1},  // diagonal \: check both diagonals
		{1, -1}, // diagonal /: check both diagonals
	}

	// Check each of the 4 possible win directions
	for _, g := range grid {
		count := 1 // Start with 1 for the piece just placed

		// Count in positive direction
		count += countDirection(cols, player, col, row, g[0], g[1])

		// Count in negative direction (opposite)
		count += countDirection(cols, player, col, row, -g[0], -g[1])

		// If 4 or more in a row found, player wins
		if count >= 4 {
			return true
		}
	}
	return false
}

// countDirection counts consecutive pieces of a player in a specific direction.
func countDirection(cols [][]int, player int, col int, row int, dc int, dr int) int {
	c := col + dc
	r := row + dr
	count := 0

	// Continue while within bounds and finding matching pieces
	for c >= 0 && c < models.Cols && r >= 0 && r < models.Rows && cols[c][r] == player {
		count += 1
		c += dc
		r += dr
	}
	return count
}
