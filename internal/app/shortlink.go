package app

import (
	"context"
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
	model.ShortURL = fmt.Sprintf("%s/%s", strings.TrimPrefix(s.cfg.BaseURL, "/"), model.Code)
	uuid, err := s.storage.Add(model)
	if err != nil {
		return model.ShortURL, err
	}

	return uuid, nil
}

func (s Service) UUID() string {
	return s.storage.UUID()
}

func (s Service) Ping(ctx context.Context) error {
	return s.storage.Ping(ctx)
}
