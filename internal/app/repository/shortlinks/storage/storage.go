package storage

import (
	"errors"
	"github.com/Orendev/shortener/internal/app/repository/shortlinks/model"
	"github.com/Orendev/shortener/internal/configs"
)

type Storage interface {
	Get(code string) (model.ShortLink, error)
	Add(shortLink model.ShortLink) (bool, error)
}

type storage struct {
	data map[string]model.ShortLink
	cfg  *configs.Configs
}

func (s *storage) Get(code string) (shortLink model.ShortLink, err error) {
	shortLink, ok := s.data[code]
	if !ok {
		err = errors.New("not found")
	}
	return
}

func (s *storage) Add(shortLink model.ShortLink) (bool, error) {
	s.data[shortLink.Code] = shortLink
	_, ok := s.data[shortLink.Code]
	return ok, nil
}

func New(cfg *configs.Configs) (Storage, error) {
	return &storage{
		cfg:  cfg,
		data: cfg.Memory,
	}, nil
}
