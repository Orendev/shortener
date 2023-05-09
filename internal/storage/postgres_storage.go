package storage

import (
	"context"
	"database/sql"
	"github.com/Orendev/shortener/internal/models"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(dsn string) (*PostgresStorage, error) {
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &PostgresStorage{
		db: db,
	}, nil
}

func (s *PostgresStorage) Close() error {
	return s.db.Close()
}

func (s PostgresStorage) GetByCode(code string) (*models.ShortLink, error) {
	return &models.ShortLink{Code: code}, nil
}

func (s PostgresStorage) Add(model *models.ShortLink) (uuid string, err error) {
	uuid = model.UUID
	err = nil
	return
}

func (s PostgresStorage) UUID() string {
	return uuid.New().String()
}

func (s *PostgresStorage) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}
