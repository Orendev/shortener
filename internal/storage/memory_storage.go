package storage

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Orendev/shortener/internal/configs"
	"github.com/Orendev/shortener/internal/models"
	"github.com/google/uuid"
)

type MemoryStorage struct {
	data map[string]models.ShortLink
	cfg  *configs.Configs
	file *File
	db   *sql.DB
}

func (s *MemoryStorage) GetByCode(code string) (*models.ShortLink, error) {
	shortLink, ok := s.data[code]
	if !ok {
		err := errors.New("not found")
		return nil, err
	}
	return &shortLink, nil
}

func (s *MemoryStorage) Add(model *models.ShortLink) (string, error) {

	s.data[model.Code] = *model

	err := s.file.Save(*model)
	if err != nil {
		return model.ShortURL, err
	}

	return model.UUID, nil
}

func (s MemoryStorage) UUID() string {
	return uuid.New().String()
}

func NewMemoryStorage(cfg *configs.Configs, db *sql.DB, file *File) (*MemoryStorage, error) {
	return &MemoryStorage{
		cfg:  cfg,
		data: cfg.Memory,
		file: file,
		db:   db,
	}, nil
}

func (s MemoryStorage) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}
