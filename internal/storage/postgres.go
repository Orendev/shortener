package storage

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
)

func NewPostgres(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return db, nil
}
