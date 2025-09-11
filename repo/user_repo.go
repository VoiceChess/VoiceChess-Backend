package repo

import (
	"database/sql"

	"samsungvoicebe/pg_sql"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateUser(userID string) error {
	_, err := r.db.Exec(pg_sql.CreateUser, userID)
	if err != nil {
		return err
	}
	return nil
}
