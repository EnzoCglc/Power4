package controllers

import (
	"errors"
	"log"
	"power4/models"
)

// play executes a complete move for a player in the specified column.
func play(game *models.GridPage, col int) error {
	// First validate that this is a legal move
	if err := validateMove(game, col); err != nil {
		return err
	}

	// Find the lowest available row in this column (gravity effect)
	row := FindAvailableRow(game.Columns, col)
	if row == -1 {
		return errors.New("colonne pleine")
	}

	// Place the piece on the board
	setPiece(game, col, row)

	// Check for win condition or switch players
	handlePostMove(game, col, row)
	return nil
}

// validateMove checks if the requested move is legal.
func validateMove(game *models.GridPage, col int) error {
	if game.GameOver {
		return errors.New("partie terminée")
	}
	if col < 0 || col >= models.Cols {
		return errors.New("colonne invalide")
	}
	return nil
}

// FindAvailableRow finds the lowest empty row in a specific column.
func FindAvailableRow(cols [][]int, col int) int {
	for row := models.Rows - 1; row >= 0; row-- {
		if cols[col][row] == models.Empty {
			return row
		}
	}
	return -1
}

// setPiece places the current player's piece at the specified position.
func setPiece(game *models.GridPage, col, row int) {
	game.Columns[col][row] = game.CurrentTurn
}

// handlePostMove processes the game state after a piece is placed.
func handlePostMove(game *models.GridPage, col, row int) {
	player := game.CurrentTurn

	// Check if this move won the game
	if VerifWin(game.Columns, player, col, row) {
		setWinner(game, player)
		return
	}

	// Check if the board is full (draw)
	if GridFull(game.Columns) {
		setDraw(game)
		return
	}

	// No win or draw, continue to next player
	switchPlayer(game)
}

// setWinner marks the game as over with a winner.
func setWinner(game *models.GridPage, player int) {
	game.Winner = player
	game.GameOver = true
	log.Printf("Player %d wins!", player)
}

// setDraw marks the game as over with no winner (draw).
func setDraw(game *models.GridPage) {
	game.GameOver = true
	game.IsDraw = true
	log.Printf("Game ended in a draw!")
}

// reset clears the game board and resets all game state to initial values.
func reset(game *models.GridPage) {
	// Clear the entire board
	for i := 0; i < models.Cols; i++ {
		game.Columns[i] = make([]int, models.Rows)
	}

	// Reset game state to defaults
	game.CurrentTurn = models.P1
	game.Winner = models.Empty
	game.GameOver = false
	game.IsDraw = false
}

// switchPlayer toggles between Player 1 and Player 2.
func switchPlayer(game *models.GridPage) {
	if game.CurrentTurn == models.P1 {
		game.CurrentTurn = models.P2
	} else {
		game.CurrentTurn = models.P1
	}
}

// GridFull checks if all columns are filled (no more moves possible).
func GridFull(cols [][]int) bool {
	for col := 0; col < models.Cols; col++ {
		if cols[col][0] == models.Empty {
			return false
		}
	}
	return true
}
