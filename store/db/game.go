package db

import (
	"database/sql"
	"log"
)

func (d *DB) JoinGame(username string) error {
	_, err := d.db.Exec("INSERT INTO Queue (username) VALUES ($1)", username)
	if err != nil {
		return err
	}
	return nil
}

func (d *DB) MatchPlayers(username string) ([]string, error) {
	// Check if there is at least one other player in the queue
	row := d.db.QueryRow("SELECT username FROM Queue WHERE username != $1 ORDER BY join_time LIMIT 1", username)
	var otherPlayerUsername string
	err := row.Scan(&otherPlayerUsername)
	if err == sql.ErrNoRows {
		// No other players in the queue
		return nil, nil
	} else if err != nil {
		// Error occurred while fetching other player
		return nil, err
	}

	// Remove the matched players from the Queue table
	_, err = d.db.Exec("DELETE FROM Queue WHERE username IN ($1, $2)", username, otherPlayerUsername)
	if err != nil {
		log.Println("Failed to remove players from queue:", err)
		return nil, err
	}

	return []string{otherPlayerUsername}, nil
}
