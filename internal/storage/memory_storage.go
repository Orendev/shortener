package storage

import (
	"errors"
	"github.com/Orendev/shortener/internal/configs"
	models "github.com/Orendev/shortener/internal/models/shortlink"
)

type MemoryStorage struct {
	data map[string]models.ShortLink
	cfg  *configs.Configs
	db   *FileDB
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

	err := s.db.Save(models.FileDB{
		OriginalURL: model.URL,
		ShortURL:    model.Code,
		UUID:        s.db.ID(),
	})
	if err != nil {
		return model.Code, err
	}

	return model.Code, nil
}

func NewMemoryStorage(cfg *configs.Configs, db *FileDB) (*MemoryStorage, error) {
	return &MemoryStorage{
		cfg:  cfg,
		data: cfg.Memory,
		db:   db,
	}, nil
}
