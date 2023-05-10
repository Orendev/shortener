package storage

import (
	"context"
	"github.com/Orendev/shortener/internal/models"
)

type ShortLinkStorage interface {
	GetByCode(ctx context.Context, code string) (*models.ShortLink, error)
	Add(ctx context.Context, model *models.ShortLink) (string, error)
	Ping(ctx context.Context) error
	Close() error
}
