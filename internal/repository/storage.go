package repository

import (
	"context"
	"errors"

	"github.com/Orendev/shortener/internal/models"
)

// ErrConflict indicates a data conflict in the storage.
var ErrConflict = errors.New("data conflict")

// ErrNotFound object not found.
var ErrNotFound = errors.New("not found")

// Storage interface for link data storage.
type Storage interface {
	GetByCode(ctx context.Context, code string) (*models.ShortLink, error)
	GetByID(ctx context.Context, id string) (*models.ShortLink, error)
	ShortLinksByUserID(ctx context.Context, userID string, limit int) ([]models.ShortLink, error)
	GetByOriginalURL(ctx context.Context, originalURL string) (*models.ShortLink, error)
	UsersStats(ctx context.Context) (int, error)
	UrlsStats(ctx context.Context) (int, error)
	Save(ctx context.Context, model models.ShortLink) error
	InsertBatch(ctx context.Context, models []models.ShortLink) error
	UpdateBatch(ctx context.Context, models []models.ShortLink) error
	DeleteFlagBatch(ctx context.Context, codes []string, userID string) error
	Ping(ctx context.Context) error
	Close() error
}
