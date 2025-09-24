package controllers

import (
	"log"
	"net/http"
	m "power4/models"
	"strconv"
)

func SwitchPlay(w http.ResponseWriter, r *http.Request) {
	action := r.FormValue("play")
	log.Println(action)

	colStr := r.FormValue("col")
	if colStr != "" {
		col, err := strconv.Atoi(colStr)
		if err != nil {
			log.Printf("invalid column %q: %v", colStr, err)
		} else {
			play(m.CurrentGame, col)
		}
	}
	render(w, "gameBoard.html", m.CurrentGame)
}

func play(game *m.GridPage, col int) {
	player := game.CurrenctTurn

	for row := m.Rows - 1; row >= 0; row-- {
		if game.Columns[col][row] == m.Empty {
			game.Columns[col][row] = player
			if verifWin(game.Columns, player, col, row) {
				log.Printf("Player %d win", player)
				return
			} else {
				if player == m.P1 {
					game.CurrenctTurn = m.P2
				} else {
					game.CurrenctTurn = m.P1
				}
				return
			}
		}
	}
}

func verifWin(cols [][]int, player int, col int, row int) bool {
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
	for c >= 0 && c < m.Cols && r >= 0 && r < m.Rows && cols[c][r] == player {
		count += 1
		c += dc
		r += dr
	}
	return count
}
