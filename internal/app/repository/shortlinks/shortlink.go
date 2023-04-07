package shortlinks

import (
	"errors"
	"fmt"
	"github.com/Orendev/shortener/internal/app/repository/shortlinks/model"
	"github.com/Orendev/shortener/internal/app/repository/shortlinks/storage"
	"github.com/Orendev/shortener/internal/configs"
	"strings"
)

type Service interface {
	Get(code string) (model.ShortLink, error)
	Add(shortLink model.ShortLink) (string, error)
}

type shortLink struct {
	storage storage.Storage
	cfg     *configs.Configs
}

func (s *shortLink) Get(code string) (shortLink model.ShortLink, err error) {
	return s.storage.Get(code)
}

func (s *shortLink) Add(shortLink model.ShortLink) (string, error) {
	ok, err := s.storage.Add(shortLink)
	if !ok {
		err = errors.New("error add")
	}

	url := fmt.Sprintf("%s/%s", strings.TrimPrefix(s.cfg.BaseURL, "/"), shortLink.Code)
	return url, err
}

func New(s storage.Storage, cfg *configs.Configs) (Service, error) {
	return &shortLink{
		storage: s,
		cfg:     cfg,
	}, nil
}
