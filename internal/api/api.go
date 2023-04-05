package api

import "github.com/Orendev/shortener/internal/app/repository/shortlinks"

type API struct {
	shortLinks *shortlinks.ShortLinks
}

func New(sl *shortlinks.ShortLinks) (*API, error) {
	return &API{
		shortLinks: sl,
	}, nil
}
