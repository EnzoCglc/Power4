package controllers

import (
	"log"
	"net/http"
	"power4/models"
	"power4/utils"
	"strconv"
)

func GameDuo(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		reset(models.CurrentGame)
		models.CurrentGame.GameMode = "duo"
	}
	utils.Render(w, "gameBoard.html", models.CurrentGame)
}

func SwitchPlay(w http.ResponseWriter, r *http.Request) {
	exit := r.FormValue("exit")

	if exit == "reset" {
		reset(models.CurrentGame)
		utils.Render(w, "index.html", nil)
		return
	}

	colStr := r.FormValue("col")
	if colStr != "" {
		col, err := strconv.Atoi(colStr)
		if err != nil {
			log.Printf("invalid column %q: %v", colStr, err)
		} else {
			play(models.CurrentGame, col)
		}
	}
	utils.Render(w, "gameBoard.html", models.CurrentGame)
}

func play(game *models.GridPage, col int) {
	player := game.CurrenctTurn

	for row := models.Rows - 1; row >= 0; row-- {
		if game.Columns[col][row] == models.Empty {
			game.Columns[col][row] = player
			if verifWin(game.Columns, player, col, row) {
				log.Printf("Player %d win", player)
				return
			} else {
				if player == models.P1 {
					game.CurrenctTurn = models.P2
				} else {
					game.CurrenctTurn = models.P1
				}
				return
			}
		}
	}
}

func reset(game *models.GridPage) {
	for i := 0; i < models.Cols; i++ {
		game.Columns[i] = make([]int, models.Rows)
	}
	game.CurrenctTurn = models.P1
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
	for c >= 0 && c < models.Cols && r >= 0 && r < models.Rows && cols[c][r] == player {
		count += 1
		c += dc
		r += dr
	}
	return count
}
