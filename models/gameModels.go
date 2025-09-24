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
	CurrenctTurn int
}

var CurrentGame = newGrid()

func newGrid() *GridPage {
	g := &GridPage{
		Columns:      make([][]int, Cols),
		CurrenctTurn: P1,
	}
	for i := 0; i < Cols; i++ {
		g.Columns[i] = make([]int, Rows)
	}
	return g
}
