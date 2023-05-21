package storage

import (
	"context"
	"errors"
	"github.com/Orendev/shortener/internal/models"
)

// ErrConflict указывает на конфликт данных в хранилище
var ErrConflict = errors.New("data conflict")

type ShortLinkStorage interface {
	GetByCode(ctx context.Context, code string) (*models.ShortLink, error)
	GetByID(ctx context.Context, id string) (*models.ShortLink, error)
	ShortLinksByUserId(ctx context.Context, userID string, limit int) ([]models.ShortLink, error)
	GetByOriginalURL(ctx context.Context, originalURL string) (*models.ShortLink, error)
	Save(ctx context.Context, model models.ShortLink) error
	InsertBatch(ctx context.Context, models []models.ShortLink) error
	UpdateBatch(ctx context.Context, models []models.ShortLink) error
	Ping(ctx context.Context) error
	Close() error
}
