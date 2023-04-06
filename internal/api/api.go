package api

import (
	"github.com/Orendev/shortener/internal/app/repository/shortlinks"
	"github.com/Orendev/shortener/internal/configs"
)

type API struct {
	shortLinks *shortlinks.ShortLinks
	cfg        *configs.Configs
}

func New(cfg *configs.Configs, sl *shortlinks.ShortLinks) (*API, error) {
	return &API{
		shortLinks: sl,
		cfg:        cfg,
	}, nil
}
