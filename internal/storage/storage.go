package storage

import (
	"context"
	"github.com/Orendev/shortener/internal/models"
)

type ShortLinkStorage interface {
	GetByCode(ctx context.Context, code string) (*models.ShortLink, error)
	GetById(ctx context.Context, id string) (*models.ShortLink, error)
	Save(ctx context.Context, model models.ShortLink) error
	InsertBatch(ctx context.Context, models []models.ShortLink) error
	UpdateBatch(ctx context.Context, models []models.ShortLink) error
	Ping(ctx context.Context) error
	Close() error
}
