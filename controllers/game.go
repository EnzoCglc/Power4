package controllers

import (
	"log"
	"net/http"
)

// Create 2D matrix to represents the game board
type GridPage struct {
	Columns [][]int
}

func SwitchPlay(w http.ResponseWriter, r *http.Request) {
	//create 7 columns
	cols := make([][]int, 7)
	for i := 0; i < 7; i++ {
		cols[i] = make([]int, 6) // create 6 case empty
	}
	if r.Method == "POST" {
		// Detect start game
		name := r.FormValue("play")
		log.Println(name)
		//Detect clic in column
		col := r.FormValue("col")
		play(col)

		data := GridPage{Columns: cols}
		render(w, "gameBoard.html", data)
	}
}

func play(col string) {
	log.Println("col", col, "has clicked")
}
