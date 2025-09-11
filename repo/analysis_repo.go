package repo

import (
	"database/sql"

	"samsungvoicebe/models"
	"samsungvoicebe/pg_sql"
)

type AnalysisRepo struct {
	db *sql.DB
}

func NewAnalysisRepo(db *sql.DB) *AnalysisRepo {
	return &AnalysisRepo{db: db}
}

func (r *AnalysisRepo) GetMoveByOrder(moveOrder int, gameID string) (models.Move, error) {
	var move models.Move
	err := r.db.QueryRow(pg_sql.GetMoveByOrder, gameID, moveOrder).Scan(&move.Move, &move.Fen)
	if err != nil {
		return models.Move{}, err
	}
	return move, nil
}

func (r *AnalysisRepo) GetGameHistoryList(userID string) ([]models.Game, error) {
	var games []models.Game
	rows, err := r.db.Query(pg_sql.GetGameHistoryList, userID)
	if err != nil {
		return []models.Game{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var game models.Game
		if err := rows.Scan(&game.GameID, &game.Date, &game.MoveAmount); err != nil {
			return []models.Game{}, err
		}
		games = append(games, game)
	}

	if err := rows.Err(); err != nil {
		return []models.Game{}, err
	}

	return games, nil
}
