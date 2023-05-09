package storage

import (
	"context"
	"errors"
	"github.com/Orendev/shortener/internal/config"
	"github.com/Orendev/shortener/internal/models"
	"github.com/google/uuid"
)

type MemoryStorage struct {
	data map[string]models.ShortLink
	cfg  *config.Configs
	file *File
	db   *PostgresStorage
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

func NewMemoryStorage(cfg *config.Configs, db *PostgresStorage, file *File) (*MemoryStorage, error) {
	return &MemoryStorage{
		cfg:  cfg,
		data: cfg.Memory,
		file: file,
		db:   db,
	}, nil
}

func (s MemoryStorage) Ping(ctx context.Context) error {
	return s.db.Ping(ctx)
}
