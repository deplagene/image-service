package db

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewPostgresStorage(connUrl string) (*sql.DB, error) {
	db, err := sql.Open("pgx", connUrl)
	if err != nil {
		return nil, err
	}

	return db, nil
}
