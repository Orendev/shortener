package memory

import (
	"context"

	"github.com/Orendev/shortener/internal/models"
	"github.com/Orendev/shortener/internal/repository"
)

// Memory - structure describing the Memory.
type Memory struct {
	data map[string]models.ShortLink
	file *File
}

// NewMemory - constructor a new instance of Memory.
func NewMemory(filePath string) (*Memory, error) {

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

// GetByCode we get a model models.ShortLink of a short link by code.
func (s *Memory) GetByCode(_ context.Context, code string) (*models.ShortLink, error) {
	shortLink, ok := s.data[code]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return &shortLink, nil
}

// GetByID we get a model models.ShortLink of a short link by id.
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
		return nil, repository.ErrNotFound
	}

	return &shortLink, nil
}

// ShortLinksByUserID we will get a list of the user's short link models.ShortLink.
func (s *Memory) ShortLinksByUserID(_ context.Context, userID string, limit int) ([]models.ShortLink, error) {
	shortLinks := make([]models.ShortLink, 0, limit)

	for _, link := range s.data {
		if link.UserID == userID {
			shortLinks = append(shortLinks, link)
		}
	}

	return shortLinks, nil
}

// GetByOriginalURL we will get the model with a short link models.ShortLink to the original URL.
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
		return nil, repository.ErrNotFound
	}

	return &shortLink, nil
}

// Save let's save the model of the short link models.ShortLink.
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

// InsertBatch group insertion of short link models []models.ShortLink.
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

// UpdateBatch group update of short link models []models.ShortLink.
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

// DeleteFlagBatch group delete of short link models []models.ShortLink.
func (s *Memory) DeleteFlagBatch(ctx context.Context, codes []string, userID string) error {
	for _, code := range codes {
		model, err := s.GetByCode(ctx, code)
		if err != nil {
			continue
		}
		if model.UserID == userID {
			model.DeletedFlag = true
			s.data[code] = *model
		}
	}
	err := s.file.Save(s.data)
	if err != nil {
		return err
	}
	return nil

}

// Ping service check.
func (s *Memory) Ping(_ context.Context) error {
	return nil
}

// Close closing the service.
func (s *Memory) Close() error {
	return nil
}
