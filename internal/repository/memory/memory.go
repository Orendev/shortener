package memory

import (
	"context"
	"errors"

	"github.com/Orendev/shortener/internal/models"
	"github.com/Orendev/shortener/internal/repository"
)

type Memory struct {
	data map[string]models.ShortLink
	file *File
}

func NewRepository(filePath string) (*Memory, error) {

	file := NewFile(filePath)

	data, err := file.Data()
	if err != nil {
		return nil, err
	}

	return &Memory{
		data: data,
		file: file,
	}, nil
}

func (s *Memory) GetByCode(_ context.Context, code string) (*models.ShortLink, error) {
	shortLink, ok := s.data[code]
	if !ok {
		err := errors.New("not found")
		return nil, err
	}
	return &shortLink, nil
}

func (s *Memory) GetByID(_ context.Context, id string) (*models.ShortLink, error) {
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

func (s *Memory) ShortLinksByUserID(_ context.Context, userID string, limit int) ([]models.ShortLink, error) {
	shortLinks := make([]models.ShortLink, 0, limit)

	for _, link := range s.data {
		if link.UserID == userID {
			shortLinks = append(shortLinks, link)
		}
	}

	return shortLinks, nil
}

func (s *Memory) GetByOriginalURL(_ context.Context, originalURL string) (*models.ShortLink, error) {
	var shortLink models.ShortLink
	ok := false

	for _, link := range s.data {

		if link.OriginalURL == originalURL {
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

func (s *Memory) Save(_ context.Context, model models.ShortLink) error {

	for _, link := range s.data {

		if link.OriginalURL == model.OriginalURL {
			return repository.ErrConflict
		}
	}
	model.DeletedFlag = false
	s.data[model.Code] = model

	err := s.file.Save(s.data)
	if err != nil {
		return err
	}

	return nil
}

func (s *Memory) InsertBatch(_ context.Context, shortLinks []models.ShortLink) error {

	for _, link := range shortLinks {
		link.DeletedFlag = false
		s.data[link.Code] = link
	}
	err := s.file.Save(s.data)
	if err != nil {
		return err
	}

	return nil
}

func (s *Memory) UpdateBatch(_ context.Context, shortLinks []models.ShortLink) error {
	for _, link := range shortLinks {
		s.data[link.Code] = link
	}
	err := s.file.Save(s.data)
	if err != nil {
		return err
	}
	return nil
}

func (s *Memory) Close() error {
	return nil
}

func (s *Memory) Ping(_ context.Context) error {
	return nil
}
