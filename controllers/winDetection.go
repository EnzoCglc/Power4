package controllers

import (
	"power4/models"
)

func VerifWin(cols [][]int, player int, col int, row int) bool {
	grid := [][2]int{
		{1, 0},  // horizontal
		{0, 1},  // vertical
		{1, 1},  // diagonal \
		{1, -1}, //diagonal /
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
