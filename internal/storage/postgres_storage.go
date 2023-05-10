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

func NewPostgresStorage(ctx context.Context, dsn string) (*PostgresStorage, error) {

	db, err := sql.Open("pgx", dsn)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	err = createTable(ctx, db)

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

func (s PostgresStorage) GetByCode(ctx context.Context, code string) (*models.ShortLink, error) {

	model := models.ShortLink{}

	// делаем запрос
	sqlStatement := `SELECT id, code, short_url, original_url  FROM short_links WHERE code = $1 LIMIT 1`
	row := s.db.QueryRowContext(ctx,
		sqlStatement, code)

	// разбираем результат
	err := row.Scan(&model.UUID, &model.Code, &model.ShortURL, &model.OriginalURL)
	if err != nil {
		return nil, err
	}

	return &model, nil
}

func (s PostgresStorage) Add(ctx context.Context, model *models.ShortLink) (string, error) {
	sqlStatement := `
	INSERT INTO short_links (id, code, short_url, original_url)
	VALUES ($1, $2, $3, $4)`

	_, err := s.db.ExecContext(
		ctx,
		sqlStatement, model.UUID, model.Code, model.ShortURL, model.OriginalURL,
	)

	if err != nil {
		return "", err
	}

	return model.UUID, nil
}

func (s PostgresStorage) UUID() string {
	return uuid.New().String()
}

func (s *PostgresStorage) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

func createTable(ctx context.Context, db *sql.DB) error {
	sqlStatement := `
	CREATE TABLE IF NOT EXISTS short_links (
	    id UUID NOT NULL primary key, 
	    code VARCHAR(255) NOT NULL UNIQUE, 
	    short_url TEXT NOT NULL, 
	    original_url TEXT NOT NULL
	    )`

	_, err := db.ExecContext(
		ctx,
		sqlStatement,
	)

	return err
}
