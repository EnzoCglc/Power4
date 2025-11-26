package models

import (
	"database/sql"
	"log"
)

// Database holds the SQLite database connection and metadata.
type Database struct {
	Connect  *sql.DB // Active database connection
	Datafile string  // Path to the SQLite database file
}

// User represents a player account in the system.
type User struct {
	ID           int    // Unique user identifier (primary key)
	Username     string // Unique username for login and display
	PasswordHash string // Bcrypt-hashed password for authentication
	Elo          int    // Current ELO rating (default: 1000)
	Win          int    // Total number of wins
	Losses       int    // Total number of losses
	CreatedAt    string // Account creation timestamp
}

// DB is the global database connection instance.
var DB *Database

// GetUserByUsername retrieves a user account by username from the database.
func GetUserByUsername(username string) (*User, error) {
	// Verify database connection exists
	if DB == nil || DB.Connect == nil {
		return nil, sql.ErrConnDone
	}

	query := `SELECT id, username, password_hash, elo, victoires, defaites, created_at
	          FROM users WHERE username = ?`

	var user User
	err := DB.Connect.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.PasswordHash,
		&user.Elo, &user.Win, &user.Losses, &user.CreatedAt,
	)

	// User not found is not an error - return nil user
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UserExists checks if a username is already registered in the database.
func UserExists(username string) (bool, error) {
	// Verify database connection exists
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

// CreateUser inserts a new user account into the database with default values.
func CreateUser(username, passwordHash string) error {
	// Verify database connection exists
	if DB == nil || DB.Connect == nil {
		return sql.ErrConnDone
	}

	query := `INSERT INTO users (username, password_hash, elo, victoires, defaites)
	          VALUES (?, ?, 1000, 0, 0)`

	_, err := DB.Connect.Exec(query, username, passwordHash)
	return err
}

// UpdateUserEloAndStats updates a user's ELO rating and win/loss statistics.
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

// UpdatePassword changes a user's password in the database.
func UpdatePassword(username, newPasswordHash string) error {
	if DB == nil || DB.Connect == nil {
		return sql.ErrConnDone
	}

	log.Printf("[DB] Updating password for user '%s'\n", username)

	result, err := execPasswordUpdate(username, newPasswordHash)
	if err != nil {
		return err
	}

	return verifyPasswordUpdate(username, result)
}

// execPasswordUpdate executes the password update query.
func execPasswordUpdate(username, newPasswordHash string) (sql.Result, error) {
	query := `UPDATE users SET password_hash = ? WHERE username = ?`
	result, err := DB.Connect.Exec(query, newPasswordHash, username)
	if err != nil {
		log.Printf("[DB] ❌ Password update failed for user '%s': %v\n", username, err)
		return nil, err
	}
	return result, nil
}

// verifyPasswordUpdate verifies that the password update affected a row.
func verifyPasswordUpdate(username string, result sql.Result) error {
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		log.Printf("[DB] ⚠️ No rows affected for user '%s'\n", username)
		return sql.ErrNoRows
	}

	log.Printf("[DB] ✅ Password updated successfully for user '%s'\n", username)
	return nil
}
