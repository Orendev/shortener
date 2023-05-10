package services

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

func (s *Service) GetByCode(ctx context.Context, code string) (*models.ShortLink, error) {
	return s.storage.GetByCode(ctx, code)
}

func (s *Service) Add(ctx context.Context, model *models.ShortLink) (string, error) {
	uuid, err := s.storage.Add(ctx, model)
	if err != nil {
		return model.ShortURL, err
	}

	return uuid, nil
}

func (s Service) Ping(ctx context.Context) error {
	return s.storage.Ping(ctx)
}

func (s Service) Close() error {
	return s.storage.Close()
}
