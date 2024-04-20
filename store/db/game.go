package db

import (
	"database/sql"
	"fmt"
	"log"
)

type Status struct {
	StatusID  int    `json:"status_id`
	MatchID   int    `json:"match_id"`
	Username  string `json:"username"`
	Opponent  string `json:"opponent"`
	IsPlaying bool   `json:"isplaying"`
}

func (d *DB) JoinGame(username string) error {
	_, err := d.db.Exec("INSERT INTO Queue (username) VALUES ($1)", username)
	if err != nil {
		return err
	}
	return nil
}

func (d *DB) MatchPlayers(username string) ([]string, error) {
	// Start a transaction
	tx, err := d.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		// Rollback the transaction if there's an error
		if err != nil {
			tx.Rollback()
		}
	}()

	// Check if there is at least one other player in the queue
	var otherPlayerUsername string
	err = tx.QueryRow("SELECT username FROM Queue WHERE username != $1 ORDER BY join_time LIMIT 1", username).Scan(&otherPlayerUsername)
	if err == sql.ErrNoRows {
		// No other players in the queue
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	// Remove the matched players from the Queue table
	_, err = tx.Exec("DELETE FROM Queue WHERE username IN ($1, $2)", username, otherPlayerUsername)
	if err != nil {
		log.Println("Failed to remove players from queue:", err)
		return nil, err
	}

	// Insert a new match record
	var matchID int
	err = tx.QueryRow("INSERT INTO Match (white_player_username, black_player_username) VALUES ($1, $2) RETURNING match_id", username, otherPlayerUsername).Scan(&matchID)
	if err != nil {
		return nil, err
	}

	// Check if the user is already in the status table
	var existingUser string
	err = tx.QueryRow("SELECT username FROM status WHERE username = $1", username).Scan(&existingUser)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if existingUser != "" {
		// If the user already exists in the status table, update their status
		_, err = tx.Exec("UPDATE status SET match_id = $1 WHERE username = $2", matchID, username)
		if err != nil {
			return nil, err
		}
		_, err = tx.Exec("UPDATE status SET match_id = $1 WHERE username = $2", matchID, otherPlayerUsername)
		if err != nil {
			return nil, err
		}
	} else {
		// If the user doesn't exist in the status table, insert a new record
		_, err = tx.Exec("INSERT INTO status (match_id, username, opponent, isPlaying) VALUES ($1, $2, $3, true), ($1, $3, $2, true)", matchID, username, otherPlayerUsername)
		if err != nil {
			return nil, err
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return []string{otherPlayerUsername}, nil
}

func (d *DB) GetUserStatus(username string) ([]Status, error) {
	// Query the status table based on the provided username
	fmt.Println("finfinag status now", username)
	rows, err := d.db.Query("SELECT * FROM status WHERE username = $1", username)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// Iterate over the rows and parse the data into Status structs
	var statuses []Status
	for rows.Next() {
		var status Status
		err := rows.Scan(&status.StatusID, &status.MatchID, &status.Username, &status.Opponent, &status.IsPlaying)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		statuses = append(statuses, status)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// If no data is found, return an empty slice without error
	if len(statuses) == 0 {
		return make([]Status, 0), nil
	}

	return statuses, nil
}
