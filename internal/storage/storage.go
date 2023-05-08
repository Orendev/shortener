package storage

import (
	"github.com/Orendev/shortener/internal/models"
)

type ShortLinkStorage interface {
	Get(code string) (*models.ShortLink, error)
	Add(model *models.ShortLink) (string, error)
}
