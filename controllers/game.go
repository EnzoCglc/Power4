package controllers

import (
	"log"
	"net/http"
)

func SwitchPlay(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("play")
		log.Println(name)
		render(w, "gameBoard.html" , nil)
	}
}
