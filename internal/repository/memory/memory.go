package memory

import (
	"context"
	"errors"

	"github.com/Orendev/shortener/internal/models"
	"github.com/Orendev/shortener/internal/repository"
)

type Repository struct {
	data map[string]models.ShortLink
	file *repository.File
}

func (s *Repository) GetByCode(_ context.Context, code string) (*models.ShortLink, error) {
	shortLink, ok := s.data[code]
	if !ok {
		err := errors.New("not found")
		return nil, err
	}
	return &shortLink, nil
}

func (s *Repository) GetByID(_ context.Context, id string) (*models.ShortLink, error) {
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

func (s *Repository) ShortLinksByUserID(_ context.Context, userID string, limit int) ([]models.ShortLink, error) {
	shortLinks := make([]models.ShortLink, 0, limit)

	for _, link := range s.data {
		if link.UserID == userID {
			shortLinks = append(shortLinks, link)
			break
		}
	}

	return shortLinks, nil
}

func (s *Repository) GetByOriginalURL(_ context.Context, originalURL string) (*models.ShortLink, error) {
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

func (s *Repository) Save(_ context.Context, model models.ShortLink) error {

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

func (s *Repository) InsertBatch(_ context.Context, shortLinks []models.ShortLink) error {

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

func (s *Repository) UpdateBatch(_ context.Context, shortLinks []models.ShortLink) error {
	for _, link := range shortLinks {
		s.data[link.Code] = link
	}
	err := s.file.Save(s.data)
	if err != nil {
		return err
	}
	return nil
}

func (s *Repository) Close() error {
	return nil
}

func NewRepository(data map[string]models.ShortLink, file *repository.File) (*Repository, error) {
	return &Repository{
		data: data,
		file: file,
	}, nil
}

func (s *Repository) Ping(_ context.Context) error {
	return nil
}
