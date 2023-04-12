package shortlink

import (
	model "github.com/Orendev/shortener/internal/app/models/shortlink"
)

type ShortLinkRepository interface {
	Get(code string) (*model.ShortLink, error)
	Add(model model.ShortLink) (string, error)
}
