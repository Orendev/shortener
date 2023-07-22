package repository

import (
	"context"
	"errors"

	"github.com/Orendev/shortener/internal/models"
)

// ErrConflict указывает на конфликт данных в хранилище
var ErrConflict = errors.New("data conflict")

type Storage interface {
	// GetByCode получим ссылки по code
	GetByCode(ctx context.Context, code string) (*models.ShortLink, error)
	// GetByID получим ссылки по id
	GetByID(ctx context.Context, id string) (*models.ShortLink, error)
	// ShortLinksByUserID получим ссылку рользователя
	ShortLinksByUserID(ctx context.Context, userID string, limit int) ([]models.ShortLink, error)
	// GetByOriginalURL получим ссылки по originalUR
	GetByOriginalURL(ctx context.Context, originalURL string) (*models.ShortLink, error)
	// Save сохраним ссылку
	Save(ctx context.Context, model models.ShortLink) error
	// InsertBatch массовая вставка ссылок
	InsertBatch(ctx context.Context, models []models.ShortLink) error
	// UpdateBatch массовое обновления
	UpdateBatch(ctx context.Context, models []models.ShortLink) error
	// Ping проверка сервиса
	Ping(ctx context.Context) error
	// Close закрываем сервис
	Close() error
}
