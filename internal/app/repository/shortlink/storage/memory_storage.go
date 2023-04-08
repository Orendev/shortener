package storage

import (
	"errors"
	"fmt"
	"github.com/Orendev/shortener/internal/app/repository/shortlink/model"
	"github.com/Orendev/shortener/internal/configs"
	"strings"
)

type MemoryStorage struct {
	data map[string]model.ShortLink
	cfg  *configs.Configs
}

func (s *MemoryStorage) Get(code string) (*model.ShortLink, error) {
	shortLink, ok := s.data[code]
	if !ok {
		err := errors.New("not found")
		return nil, err
	}
	return &shortLink, nil
}

func (s *MemoryStorage) Add(shortLink model.ShortLink) (string, error) {
	s.data[shortLink.Code] = shortLink
	url := fmt.Sprintf("%s/%s", strings.TrimPrefix(s.cfg.BaseURL, "/"), shortLink.Code)
	return url, nil
}

func NewMemoryStorage(cfg *configs.Configs) (*MemoryStorage, error) {
	return &MemoryStorage{
		cfg:  cfg,
		data: cfg.Memory,
	}, nil
}
