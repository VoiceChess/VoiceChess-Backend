package repo

import (
	"database/sql"

	"samsungvoicebe/pg_sql"
)

type GameplayRepo struct {
	db *sql.DB
}

func NewGameplayRepo(db *sql.DB) *GameplayRepo {
	return &GameplayRepo{db: db}
}

func (r *GameplayRepo) GameMove(gameID, fen, move string) error {
	_, err := r.db.Exec(pg_sql.Move, gameID, fen, move)
	if err != nil {
		return err
	}
	return nil
}

func (r *GameplayRepo) CreateGame(userID string) (string, error) {
	var gameID string
	err := r.db.QueryRow(pg_sql.CreateGame, userID).Scan(&gameID)
	if err != nil {
		return "", err
	}
	return gameID, nil
}
