package models

import (
	"log"
	"time"
)

type History struct {
	ID        int       `json:"id"`
	Player1   string    `json:"player1"`
	Player2   string    `json:"player2"`
	Winner    string    `json:"winner"`
	Delta     int       `json:"delta"`
	Ranked    bool      `json:"ranked"`
	Date      time.Time `json:"date"`
}

func GetHistoryByPlayer(username string) ([]History, error) {
	query := `
		SELECT id, player1, player2, winner, delta, ranked, date
		FROM match_history
		WHERE player1 = ? OR player2 = ?
		ORDER BY date DESC;
	`
	rows, err := DB.Connect.Query(query, username, username)
	if err != nil {
		log.Println("❌ Failed to query match_history table:", err)
		return nil, err
	}
	defer rows.Close()

	var history []History
	for rows.Next() {
		var h History
		err := rows.Scan(
			&h.ID,
			&h.Player1,
			&h.Player2,
			&h.Winner,
			&h.Delta,
			&h.Ranked,
			&h.Date,
		)
		if err != nil {
			log.Println("⚠️ Failed to scan row:", err)
			continue
		}
		history = append(history, h)
	}

	log.Printf("📜 %d match(es) retrieved for user '%s'\n", len(history), username)
	return history, nil
}


func InsertHistory(player1, player2, winner string, delta int, ranked bool) error {
	query := `
	INSERT INTO match_history (player1, player2, winner, delta, ranked)
	VALUES (?, ?, ?, ?, ?);
	`

	_, err := DB.Connect.Exec(query, player1, player2, winner, delta, ranked)
	if err != nil {
		log.Printf("❌ Failed to insert match into history: %v\n", err)
		return err
	}

	log.Printf("✅ Match inserted: %s vs %s | Winner: %s | Δ%d | ranked=%v\n",
		player1, player2, winner, delta, ranked)

	return nil
}


