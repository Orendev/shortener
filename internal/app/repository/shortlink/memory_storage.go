package shortlink

import (
	"errors"
	models "github.com/Orendev/shortener/internal/app/models/shortlink"
	"github.com/Orendev/shortener/internal/configs"
)

type MemoryStorage struct {
	data map[string]models.ShortLink
	cfg  *configs.Configs
}

func (s *MemoryStorage) Get(code string) (*models.ShortLink, error) {
	shortLink, ok := s.data[code]
	if !ok {
		err := errors.New("not found")
		return nil, err
	}
	return &shortLink, nil
}

func (s *MemoryStorage) Add(model *models.ShortLink) (string, error) {
	s.data[model.Code] = *model
	return model.Code, nil
}

func NewMemoryStorage(cfg *configs.Configs) (*MemoryStorage, error) {
	return &MemoryStorage{
		cfg:  cfg,
		data: cfg.Memory,
	}, nil
}
