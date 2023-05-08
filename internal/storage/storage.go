package storage

import (
	"context"
	"github.com/Orendev/shortener/internal/models"
)

type ShortLinkStorage interface {
	GetByCode(code string) (*models.ShortLink, error)
	Add(model *models.ShortLink) (string, error)
	UUID() string
	Ping(ctx context.Context) error
}
