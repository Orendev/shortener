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

func (s *Service) GetById(ctx context.Context, id string) (*models.ShortLink, error) {
	return s.storage.GetById(ctx, id)
}

func (s *Service) Save(ctx context.Context, model models.ShortLink) error {
	err := s.storage.Save(ctx, model)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) InsertBatch(ctx context.Context, shortLinks []models.ShortLink) error {
	return s.storage.InsertBatch(ctx, shortLinks)
}

func (s *Service) UpdateBatch(ctx context.Context, shortLinks []models.ShortLink) error {
	return s.storage.UpdateBatch(ctx, shortLinks)
}

func (s Service) Ping(ctx context.Context) error {
	return s.storage.Ping(ctx)
}

func (s Service) Close() error {
	return s.storage.Close()
}
