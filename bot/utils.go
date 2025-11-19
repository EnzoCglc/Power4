package bot

import (
	"power4/models"
	"power4/controllers"
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

func SimulateMove(game *models.GridPage, col int, player int) int {
	row := controllers.FindAvailableRow(game.Columns, col)
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

func CheckWin(game *models.GridPage, player int, col int, row int) bool {
	if row == -1 {
		return false
	}
	return controllers.VerifWin(game.Columns, player, col, row)
}

func GetNextPlayer(currentplayer int) int {
	if currentplayer == models.P1 {
		return models.P2
	}
	return models.P1
}

func IsBoardFull(game *models.GridPage) bool {
	return controllers.GridFull(game.Columns)
}