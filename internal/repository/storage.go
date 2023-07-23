package repository

import (
	"context"
	"errors"

	"github.com/Orendev/shortener/internal/models"
)

// ErrConflict указывает на конфликт данных в хранилище
var ErrConflict = errors.New("data conflict")

type Storage interface {
	// GetByCode we get a model models.ShortLink of a short link by code.
	GetByCode(ctx context.Context, code string) (*models.ShortLink, error)
	// GetByID we get a model models.ShortLink of a short link by id.
	GetByID(ctx context.Context, id string) (*models.ShortLink, error)
	// ShortLinksByUserID we will get a list of the user's short link models.ShortLink.
	ShortLinksByUserID(ctx context.Context, userID string, limit int) ([]models.ShortLink, error)
	// GetByOriginalURL we will get the model with a short link models.ShortLink to the original URL.
	GetByOriginalURL(ctx context.Context, originalURL string) (*models.ShortLink, error)
	// Save let's save the model of the short link models.ShortLink.
	Save(ctx context.Context, model models.ShortLink) error
	// InsertBatch group insertion of short link models []models.ShortLink.
	InsertBatch(ctx context.Context, models []models.ShortLink) error
	// UpdateBatch group update of short link models []models.ShortLink.
	UpdateBatch(ctx context.Context, models []models.ShortLink) error
	// Ping service check.
	Ping(ctx context.Context) error
	// Close closing the service.
	Close() error
}
