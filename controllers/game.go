package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"power4/models"
	"power4/utils"
)

func GameDuo(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		reset(models.CurrentGame)
		models.CurrentGame.GameMode = "duo"
	}
	utils.Render(w, "gameBoard.html", models.CurrentGame)
	log.Println("Duo mod active")
}

func SwitchPlay(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Col int `json:"col"`
		Exit string `json:"reset"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		JSONError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if request.Exit == "reset" {
		reset(models.CurrentGame)
		utils.Render(w, "index.html", nil)
		return
	}

	if request.Col >= 0 {
		err = play(models.CurrentGame, request.Col)
		if err != nil {
			log.Printf("Invalid move in column %d: %v", request.Col, err)
			return
		}

		log.Printf("Player %d played in column %d", models.CurrentGame.CurrenctTurn, request.Col)
	}

	JSONSuccess(w, map[string]interface{}{
		"game": models.CurrentGame,
	})
}

func JSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := map[string]interface{}{
		"success": false,
		"error":   message,
	}

	json.NewEncoder(w).Encode(response)
}

func JSONSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"success": true,
		"data":    data,
	}

	json.NewEncoder(w).Encode(response)
}