package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Orendev/shortener/internal/models"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(dsn string) (*PostgresStorage, error) {

	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}

	return &PostgresStorage{
		db: db,
	}, nil
}

func (s *PostgresStorage) Close() error {
	return s.db.Close()
}

func (s *PostgresStorage) GetByCode(ctx context.Context, code string) (*models.ShortLink, error) {

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

func (s *PostgresStorage) GetByID(ctx context.Context, id string) (*models.ShortLink, error) {

	model := models.ShortLink{}

	stmt, err := s.db.PrepareContext(ctx,
		`SELECT id, code, short_url, original_url  FROM short_links WHERE id = $1 LIMIT 1`)

	if err != nil {
		return nil, err
	}

	defer func() {
		err := stmt.Close()
		if err != nil {
			_ = fmt.Errorf("db shutdown: %w", err)
		}
	}()

	// делаем запрос
	row := stmt.QueryRowContext(ctx, id)

	// разбираем результат
	err = row.Scan(&model.UUID, &model.Code, &model.ShortURL, &model.OriginalURL)
	if err != nil {
		return nil, err
	}

	return &model, nil
}

func (s *PostgresStorage) GetByOriginalURL(ctx context.Context, originalURL string) (*models.ShortLink, error) {

	model := models.ShortLink{}

	stmt, err := s.db.PrepareContext(ctx,
		`SELECT id, code, short_url, original_url  FROM short_links WHERE original_url = $1 LIMIT 1`)

	if err != nil {
		return nil, err
	}

	defer func() {
		err := stmt.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// делаем запрос
	row := stmt.QueryRowContext(ctx, originalURL)

	// разбираем результат
	err = row.Scan(&model.UUID, &model.Code, &model.ShortURL, &model.OriginalURL)
	if err != nil {
		return nil, err
	}

	return &model, nil
}

func (s *PostgresStorage) Save(ctx context.Context, model models.ShortLink) error {
	sqlStatement := `
	INSERT INTO short_links (id, code, short_url, original_url)
	VALUES ($1, $2, $3, $4)
	`

	_, err := s.db.ExecContext(
		ctx,
		sqlStatement, model.UUID, model.Code, model.ShortURL, model.OriginalURL,
	)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgerrcode.UniqueViolation == pgErr.Code {
			err = ErrConflict
		}
		return err
	}

	return nil
}

func (s *PostgresStorage) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

func (s *PostgresStorage) InsertBatch(ctx context.Context, shortLinks []models.ShortLink) error {
	// начинаем транзакцию
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx,
		`INSERT INTO short_links (id, code, short_url, original_url)
				VALUES($1, $2, $3, $4)`)
	if err != nil {
		return err
	}

	defer func() {
		err := stmt.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	for _, sl := range shortLinks {
		_, err = stmt.ExecContext(ctx, sl.UUID, sl.Code, sl.ShortURL, sl.OriginalURL)
		if err != nil {
			// если ошибка, то откатываем изменения
			errRollback := tx.Rollback()
			if errRollback != nil {
				return errRollback
			}
			return err
		}
	}
	return tx.Commit()
}

func (s *PostgresStorage) UpdateBatch(ctx context.Context, shortLinks []models.ShortLink) error {
	// начинаем транзакцию
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx,
		`UPDATE short_links SET original_url = $1 WHERE id = $2`)

	if err != nil {
		return err
	}

	defer func() {
		err := stmt.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	for _, sl := range shortLinks {
		_, err = stmt.ExecContext(ctx, sl.OriginalURL, sl.UUID)
		if err != nil {
			// если ошибка, то откатываем изменения
			errRollback := tx.Rollback()
			if errRollback != nil {
				return errRollback
			}
			return err
		}
	}
	return tx.Commit()
}

// Bootstrap подготавливает БД к работе, создавая необходимые таблицы и индексы
func (s *PostgresStorage) Bootstrap(ctx context.Context) error {

	sqlStatement := `
	CREATE TABLE IF NOT EXISTS short_links (
	    id UUID NOT NULL primary key, 
	    code VARCHAR(255) NOT NULL UNIQUE, 
	    short_url TEXT NOT NULL, 
	    original_url TEXT NOT NULL UNIQUE
	    )`

	_, err := s.db.ExecContext(
		ctx,
		sqlStatement,
	)

	return err
}
