package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"power4/models"
	"math/rand"
)

type gameResultBody struct {
	Winner  int    `json:"winner"`
	Player1 string `json:"player1"`
	Player2 string `json:"player2"`
	IsDraw  bool   `json:"isDraw"`
}

type eloResult struct {
	Winner string         `json:"winner"`
	Delta  int            `json:"delta"`
	Elo    map[string]int `json:"elo"`
}

func GameResult(w http.ResponseWriter, r *http.Request) {
	body, err := decodeBody(r)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if !models.CurrentGame.Ranked {
		writeJSON(w, http.StatusOK, map[string]string{
			"message": "Unranked game, no ELO modification",
		})
		return
	}

	result, err := processResult(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"winner":  result.Winner,
		"delta":   result.Delta,
		"elo":     result.Elo,
	})
}

func decodeBody(r *http.Request) (*gameResultBody, error) {
	var body gameResultBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}
	return &body, nil
}

func processResult(body *gameResultBody) (*eloResult, error) {
	if body.IsDraw {
		return nil, errors.New("draw game, no ELO change")
	}

	winnerName, loserName := getResult(body)
	if winnerName == "" || loserName == "" {
		return nil, errors.New("invalid players: winner/loser missing")
	}

	winner, err := models.GetUserByUsername(winnerName)
	if err != nil || winner == nil {
		return nil, errors.New("winner not found in database")
	}

	loser, err := models.GetUserByUsername(loserName)
	if err != nil || loser == nil {
		return nil, errors.New("loser not found in database")
	}

	delta := calculateElo(winner, loser)

	if err := models.UpdateUserEloAndStats(winner); err != nil {
		return nil, err
	}
	if err := models.UpdateUserEloAndStats(loser); err != nil {
		return nil, err
	}

	return &eloResult{
		Winner: winner.Username,
		Delta:  delta,
		Elo: map[string]int{
			winner.Username: winner.Elo,
			loser.Username:  loser.Elo,
		},
	}, nil
}

func getResult(body *gameResultBody) (string, string) {
	switch body.Winner {
	case models.P1:
		return body.Player1, body.Player2
	case models.P2:
		return body.Player2, body.Player1
	default:
		return "", ""
	}
}

func calculateElo(winner, loser *models.User) int {
	delta := 10 + rand.Intn(21) 
	winner.Elo += delta
	winner.Win++

	loser.Elo -= delta
	if loser.Elo < 0 {
		loser.Elo = 0
	}
	loser.Losses++

	return delta
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}