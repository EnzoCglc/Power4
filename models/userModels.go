package models

import "database/sql"
import "log"

type Database struct {
	Connect  *sql.DB
	Datafile string
}
type User struct {
	ID           int
	Username     string
	PasswordHash string
	Elo          int
	Win          int
	Losses       int
	CreatedAt    string
}

var DB *Database

func GetUserByUsername(username string) (*User, error) {
	if DB == nil || DB.Connect == nil {
		return nil, sql.ErrConnDone
	}

	query := `SELECT id, username,password_hash, elo, victoires, defaites, created_at
	          FROM users WHERE username = ?`

	var user User
	err := DB.Connect.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.PasswordHash,
		&user.Elo, &user.Win, &user.Losses, &user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func UserExists(username string) (bool, error) {
	if DB == nil || DB.Connect == nil {
		return false, sql.ErrConnDone
	}

	query := `SELECT COUNT(*) FROM users WHERE username = ?`

	var count int
	err := DB.Connect.QueryRow(query, username).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func CreateUser(username, passwordHash string) error {
	if DB == nil || DB.Connect == nil {
		return sql.ErrConnDone
	}

	query := `INSERT INTO users (username, password_hash, elo, victoires, defaites)
	          VALUES (?, ?, 1000, 0, 0)`

	_, err := DB.Connect.Exec(query, username, passwordHash)
	return err
}

func UpdateUserEloAndStats(user *User) error {
	log.Printf("[DB] Updating ELO/stats for user '%s' | ELO=%d | Wins=%d | Losses=%d\n",
		user.Username, user.Elo, user.Win, user.Losses)

	_, err := DB.Connect.Exec(`
		UPDATE users
		SET elo = ?, victoires = ?, defaites = ?
		WHERE username = ?`,
		user.Elo, user.Win, user.Losses, user.Username,
	)
	if err != nil {
		log.Printf("[DB] ❌ SQL update failed for user '%s': %v\n", user.Username, err)
		return err
	}

	return nil
}
