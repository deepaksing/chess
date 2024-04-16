package db

import (
	"database/sql"

	"github.com/deepaksing/chess/types"
)

func (d *DB) CreateUser(user *types.UserTable) error {
	sqlCommand := `INSERT INTO "User" (username, password_hash) VALUES ($1, $2)`
	_, err := d.db.Exec(sqlCommand, user.Username, user.HashedPassword)
	if err != nil {
		return err
	}
	return nil
}

func (d *DB) FindUserByUsername(username string) (*types.UserTable, error) {
	sqlCommand := `SELECT username, password_hash from "User" WHERE username = $1`

	row := d.db.QueryRow(sqlCommand, username)

	var user types.UserTable
	err := row.Scan(&user.Username, &user.HashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
