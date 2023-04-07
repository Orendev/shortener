package storage

import (
	"errors"
	"github.com/Orendev/shortener/internal/app/repository/shortlinks/Model"
	"github.com/Orendev/shortener/internal/configs"
)

type Storage interface {
	Get(code string) (Model.ShortLink, error)
	Add(shortLink Model.ShortLink) (bool, error)
}

type storage struct {
	data map[string]Model.ShortLink
	cfg  *configs.Configs
}

func (s *storage) Get(code string) (shortLink Model.ShortLink, err error) {
	shortLink, ok := s.data[code]
	if !ok {
		err = errors.New("not found")
	}
	return
}

func (s *storage) Add(shortLink Model.ShortLink) (bool, error) {
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
