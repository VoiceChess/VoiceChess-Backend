package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func New() (*sql.DB, error) {
	databaseUrl := os.Getenv("POSTGRESQL_URL")
	if databaseUrl == "" {
		return nil, fmt.Errorf("POSTGRESQL_URL environment variable not set")
	}

	db, err := sql.Open("postgres", databaseUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}

	return db, nil
}
