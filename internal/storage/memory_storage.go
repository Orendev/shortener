package storage

import (
	"context"
	"errors"
	"github.com/Orendev/shortener/internal/config"
	"github.com/Orendev/shortener/internal/models"
)

type MemoryStorage struct {
	data map[string]models.ShortLink
	cfg  *config.Configs
	file *File
}

func (s *MemoryStorage) GetByCode(_ context.Context, code string) (*models.ShortLink, error) {
	shortLink, ok := s.data[code]
	if !ok {
		err := errors.New("not found")
		return nil, err
	}
	return &shortLink, nil
}

func (s *MemoryStorage) GetByID(_ context.Context, id string) (*models.ShortLink, error) {
	var shortLink models.ShortLink
	ok := false

	for _, link := range s.data {

		if link.UUID == id {
			shortLink = link
			ok = true
			break
		}
	}

	if !ok {
		err := errors.New("not found")
		return nil, err
	}
	return &shortLink, nil
}

func (s *MemoryStorage) Save(_ context.Context, model models.ShortLink) error {

	s.data[model.Code] = model

	err := s.file.Save(s.data)
	if err != nil {
		return err
	}

	return nil
}

func (s *MemoryStorage) InsertBatch(_ context.Context, shortLinks []models.ShortLink) error {

	for _, link := range shortLinks {
		s.data[link.Code] = link
	}
	err := s.file.Save(s.data)
	if err != nil {
		return err
	}

	return nil
}

func (s *MemoryStorage) UpdateBatch(_ context.Context, shortLinks []models.ShortLink) error {
	for _, link := range shortLinks {
		s.data[link.Code] = link
	}
	err := s.file.Save(s.data)
	if err != nil {
		return err
	}
	return nil
}

func (s MemoryStorage) Close() error {
	return nil
}

func NewMemoryStorage(cfg *config.Configs, file *File) (*MemoryStorage, error) {
	return &MemoryStorage{
		cfg:  cfg,
		data: cfg.Memory,
		file: file,
	}, nil
}

func (s MemoryStorage) Ping(_ context.Context) error {
	return nil
}
