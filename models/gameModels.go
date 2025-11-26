package models

// Game board dimensions and player constants for Connect 4
const (
	Rows = 6  // Number of rows 
	Cols = 7  // Number of columns 

	Empty = 0  // Empty cell value
	P1    = 1  // Player 1 identifier
	P2    = 2  // Player 2 identifier
)

// GridPage represents the complete state of a Connect 4 game.
type GridPage struct {
	Columns     [][]int // 2D array representing the game board [column][row]
	CurrentTurn int     // Current player's turn (P1 or P2)
	GameMode    string  // Game mode: "duo" (local multiplayer) or "bot" (vs AI)
	BotLvl      int     // Bot difficulty level (1-5), only used in bot mode
	Winner      int     // Winner of the game (P1, P2, or Empty if ongoing)
	GameOver    bool    // Whether the game has ended
	IsDraw      bool    // Whether the game ended in a draw
	Ranked      bool    // Whether this game affects ELO ratings
}

// CurrentGame is the global game state instance
var CurrentGame = newGrid()

// newGrid creates and initializes a new game board with default values.
func newGrid() *GridPage {
	g := &GridPage{
		Columns:     make([][]int, Cols),
		CurrentTurn: P1,
		Winner:      Empty,
	}

	// Initialize each column with the correct number of rows
	for i := 0; i < Cols; i++ {
		g.Columns[i] = make([]int, Rows)
	}

	return g
}
