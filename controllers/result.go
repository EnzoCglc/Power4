package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"power4/models"
	"math/rand"
)

type gameResultBody struct {
	Winner  int    `json:"winner"`  // Player ID (1 or 2) who won the game
	Player1 string `json:"player1"` // Username of player 1
	Player2 string `json:"player2"` // Username of player 2
	IsDraw  bool   `json:"isDraw"`  // Whether the game ended in a draw
}

type eloResult struct {
	Winner string         `json:"winner"` // Username of the winning player
	Delta  int            `json:"delta"`  // ELO points gained/lost
	Elo    map[string]int `json:"elo"`    // Map of usernames to their new ELO ratings
}

func GameResult(w http.ResponseWriter, r *http.Request) {
	// Decode the game result from the request body
	body, err := decodeBody(r)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Skip ELO processing for unranked games
	if !models.CurrentGame.Ranked {
		writeJSON(w, http.StatusOK, map[string]string{
			"message": "Unranked game, no ELO modification",
		})
		return
	}

	// Process the result: calculate ELO changes, update database, record history
	result, err := processResult(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send back the updated ELO information to the client
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"winner":  result.Winner,
		"delta":   result.Delta,
		"elo":     result.Elo,
	})
}

// decodeBody parses the JSON request body into a gameResultBody struct.
func decodeBody(r *http.Request) (*gameResultBody, error) {
	var body gameResultBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}
	return &body, nil
}

// processResult handles the complete workflow of processing a game result.
func processResult(body *gameResultBody) (*eloResult, error) {
	// Skip processing if the game was a draw
	if body.IsDraw {
		return nil, errors.New("draw game, no ELO change")
	}

	// Determine winner and loser based on the Winner ID
	winnerName, loserName := getResult(body)
	if winnerName == "" || loserName == "" {
		return nil, errors.New("invalid players: winner/loser missing")
	}

	// Fetch user data from database
	winner, err := models.GetUserByUsername(winnerName)
	if err != nil || winner == nil {
		log.Printf("Error fetching winner (%s): %v", winnerName, err)
		return nil, errors.New("failed to fetch winner from database")
	}

	loser, err := models.GetUserByUsername(loserName)
	if err != nil || loser == nil {
		log.Printf("Error fetching loser (%s): %v", loserName, err)
		return nil, errors.New("failed to fetch loser from database")
	}

	// Calculate ELO changes (modifies user structs in-place)
	delta := calculateElo(winner, loser)

	// Persist updated ELO and stats to database
	if err := models.UpdateUserEloAndStats(winner); err != nil {
		log.Println("Failed to update winner ELO:", err)
		return nil, errors.New("failed to update winner in database")
	}

	if err := models.UpdateUserEloAndStats(loser); err != nil {
		log.Println("Failed to update loser ELO:", err)
		return nil, errors.New("failed to update loser in database")
	}

	// Record match in history table
	if err := models.InsertHistory(winner.Username, loser.Username, winner.Username, delta, models.CurrentGame.Ranked); err != nil {
		log.Println("Failed to insert match history:", err)
		return nil, errors.New("failed to insert match history")
	}

	// Return the result containing new ELO ratings
	return &eloResult{
		Winner: winner.Username,
		Delta:  delta,
		Elo: map[string]int{
			winner.Username: winner.Elo,
			loser.Username:  loser.Elo,
		},
	}, nil
}

// getResult extracts the winner and loser usernames from the game result.
func getResult(body *gameResultBody) (string, string) {
	switch body.Winner {
	case models.P1:
		return body.Player1, body.Player2
	case models.P2:
		return body.Player2, body.Player1
	default:
		log.Printf("Unknown winner ID: %d", body.Winner)
		return "", ""
	}
}

// calculateElo computes the ELO rating changes for both players after a match.
func calculateElo(winner, loser *models.User) int {
	// Generate random ELO change between 10 and 30 points
	delta := 10 + rand.Intn(21)

	// Update winner's ELO and win count
	winner.Elo += delta
	winner.Win++

	// Update loser's ELO and loss count, ensuring ELO doesn't go negative
	loser.Elo -= delta
	if loser.Elo < 0 {
		loser.Elo = 0
	}
	loser.Losses++

	return delta
}

// writeJSON sends a JSON response to the client with the specified status code.
func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}