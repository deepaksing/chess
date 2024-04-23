package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/deepaksing/chess/store"
	"github.com/deepaksing/chess/types"
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

	chessboardState := store.NewChessBoard()
	fmt.Println(chessboardState)
	// chessboardStateJSON, err := json.Marshal(chessboardState)
	// if err != nil {
	// 	return nil, err
	// }

	// Insert a new match record
	var matchID int
	err = tx.QueryRow("INSERT INTO Match (white_player_username, black_player_username, chessboard_state) VALUES ($1, $2, $3) RETURNING match_id", username, otherPlayerUsername, chessboardState).Scan(&matchID)
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

func (d *DB) GetBoardState(match_id int) (string, error) {
	fmt.Println("fetchnig board ")
	var boardState string
	err := d.db.QueryRow("SELECT chessboard_state FROM match WHERE match_id = $1", match_id).Scan(&boardState)
	fmt.Println("board state", boardState)
	if err != nil {
		return "", err
	}
	return boardState, nil
}

func (d *DB) SaveBoardMove(Move types.MoveResp) error {
	//save the move and get move_id
	fmt.Println(Move)
	var move_id int
	err := d.db.QueryRow("INSERT INTO move (move_from, move_to, move_type, player_username) VALUES ($1, $2, $3, $4) RETURNING move_id", Move.Move_from, Move.Move_to, Move.Move_type, Move.Player_username).Scan(&move_id)

	if err != nil {
		return err
	}

	//save move_id in match table

	res := d.db.QueryRow("UPDATE match SET move_ids = array_append(move_ids, $1) WHERE match_id = $2", move_id, Move.Match_id)
	if res.Err() != nil {
		return err
	}
	//update the board_state
	//1. fetch board state
	//2. modify the board and update

	var board_state string
	res = d.db.QueryRow("SELECT chessboard_state from match WHERE match_id=$1", Move.Match_id)
	res.Scan(&board_state)

	updatedBoard, err := store.UpdateBoardState(store.ConvertStringToChessboard(board_state), Move.Move_from, Move.Move_to, Move.Move_type)
	if err != nil {
		return err
	}
	fmt.Println(board_state)
	fmt.Println(updatedBoard)

	res = d.db.QueryRow("UPDATE match SET chessboard_state = $1 WHERE match_id = $2", store.ConvertChessboardToString(updatedBoard), Move.Match_id)
	if res.Err() != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
