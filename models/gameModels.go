package models

const (
	Rows = 6
	Cols = 7

	Empty = 0
	P1    = 1
	P2    = 2
)

type GridPage struct {
	Columns      [][]int
	CurrentTurn int
	GameMode     string
	BotLvl		 int
	Winner       int
	GameOver     bool
	IsDraw       bool
	Ranked		 bool
}

var CurrentGame = newGrid()

func newGrid() *GridPage {
	g := &GridPage{
		Columns:      make([][]int, Cols),
		CurrentTurn: P1,
		Winner:       Empty,
	}
	for i := 0; i < Cols; i++ {
		g.Columns[i] = make([]int, Rows)
	}
	return g
}
