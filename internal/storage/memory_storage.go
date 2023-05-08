package storage

import (
	"errors"
	"github.com/Orendev/shortener/internal/configs"
	"github.com/Orendev/shortener/internal/models"
	"github.com/google/uuid"
)

type MemoryStorage struct {
	data map[string]models.ShortLink
	cfg  *configs.Configs
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

	return model.UUID, nil
}

func (s MemoryStorage) UUID() string {
	return uuid.New().String()
}

func NewMemoryStorage(cfg *configs.Configs) (*MemoryStorage, error) {
	return &MemoryStorage{
		cfg:  cfg,
		data: cfg.Memory,
	}, nil
}
