package storage

import (
	models "github.com/Orendev/shortener/internal/models/shortlink"
)

type ShortLinkRepository interface {
	Get(code string) (*models.ShortLink, error)
	Add(model *models.ShortLink) (string, error)
}
