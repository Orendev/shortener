package storage

import (
	"github.com/Orendev/shortener/internal/app/repository/shortlink/model"
)

type ShortLinkRepository interface {
	Get(code string) (*model.ShortLink, error)
	Add(shortLink model.ShortLink) (string, error)
}
