package app

import (
	"fmt"
	"github.com/Orendev/shortener/internal/configs"
	"github.com/Orendev/shortener/internal/models"
	"github.com/Orendev/shortener/internal/storage"
	"strings"
)

type Service struct {
	storage storage.ShortLinkStorage
	cfg     *configs.Configs
}

func NewService(storage storage.ShortLinkStorage, cfg *configs.Configs) *Service {

	return &Service{
		storage: storage,
		cfg:     cfg,
	}
}

func (s *Service) GetByCode(code string) (*models.ShortLink, error) {
	return s.storage.GetByCode(code)
}

func (s *Service) Add(model *models.ShortLink) (string, error) {
	model.ShortUrl = fmt.Sprintf("%s/%s", strings.TrimPrefix(s.cfg.BaseURL, "/"), model.ShortUrl)
	uuid, err := s.storage.Add(model)
	if err != nil {
		return model.ShortUrl, err
	}

	return uuid, nil
}

func (s Service) Uuid() string {
	return s.storage.Uuid()
}
