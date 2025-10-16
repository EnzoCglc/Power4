package controllers

import (
	"errors"
	"log"
	"power4/models"
)

func play(game *models.GridPage, col int) error {
	if err := validateMove(game, col); err != nil {
		return err
	}

	row := findAvailableRow(game.Columns, col)
	if row == -1 {
		return errors.New("colonne pleine")
	}

	setPiece(game, col, row)
	handlePostMove(game, col, row)
	return nil
}

func validateMove(game *models.GridPage, col int) error {
	if game.GameOver {
		return errors.New("partie terminée")
	}
	if col < 0 || col >= models.Cols {
		return errors.New("colonne invalide")
	}
	return nil
}

func findAvailableRow(cols [][]int, col int) int {
	for row := models.Rows - 1; row >= 0; row-- {
		if cols[col][row] == models.Empty {
			return row
		}
	}
	return -1
}

func setPiece(game *models.GridPage, col, row int) {
	game.Columns[col][row] = game.CurrenctTurn
}

func handlePostMove(game *models.GridPage, col, row int) {
	player := game.CurrenctTurn

	if verifWin(game.Columns, player, col, row) {
		setWinner(game, player)
		return
	}

	if GridFull(game.Columns) {
		setDraw(game)
		return
	}

	switchPlayer(game)
}

func setWinner(game *models.GridPage, player int) {
	game.Winner = player
	game.GameOver = true
	log.Printf("Player %d wins!", player)
}

func setDraw(game *models.GridPage) {
	game.GameOver = true
	game.IsDraw = true
	log.Printf("Game ended in a draw!")
}

func reset(game *models.GridPage) {
	for i := 0; i < models.Cols; i++ {
		game.Columns[i] = make([]int, models.Rows)
	}
	game.CurrenctTurn = models.P1
	game.Winner = models.Empty
	game.GameOver = false
	game.IsDraw = false
}

func switchPlayer(game *models.GridPage) {
	if game.CurrenctTurn == models.P1 {
		game.CurrenctTurn = models.P2
	} else {
		game.CurrenctTurn = models.P1
	}
}

func GridFull(cols [][]int) bool {
	for col := 0; col < models.Cols; col++ {
		if cols[col][0] == models.Empty {
			return false
		}
	}
	return true
}
