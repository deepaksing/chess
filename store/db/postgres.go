package db

import (
	"database/sql"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

type DB struct {
	db *sql.DB
}

var (
	connStr = "user=postgres password=postgres dbname=Chess sslmode=disable"
)

func NewDB() (*DB, error) {
	//create postgres db
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return &DB{
		db: db,
	}, nil
}

func (d *DB) Migrate() error {
	content, err := os.ReadFile("schema.sql")
	if err != nil {
		return err
	}
	queries := strings.Split(string(content), ";")

	// Execute each query
	for _, query := range queries {
		// Trim leading/trailing whitespace and skip empty queries
		trimmedQuery := strings.TrimSpace(query)
		if trimmedQuery == "" {
			continue
		}

		// Execute query
		_, err := d.db.Exec(trimmedQuery)
		if err != nil {
			return err
		}
	}
	return nil
}
