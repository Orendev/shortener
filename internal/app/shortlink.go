package app

import (
	"context"
	"github.com/Orendev/shortener/internal/models"
	"github.com/Orendev/shortener/internal/storage"
)

type Service struct {
	storage storage.ShortLinkStorage
}

func NewService(storage storage.ShortLinkStorage) *Service {

	return &Service{
		storage: storage,
	}
}

func (s *Service) GetByCode(code string) (*models.ShortLink, error) {
	return s.storage.GetByCode(code)
}

func (s *Service) Add(model *models.ShortLink) (string, error) {
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
