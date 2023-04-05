package api

import "github.com/Orendev/shortener/internal/app/repository/shortlinks"

// AddLink add a new short link
func (a *API) AddLink(sl shortlinks.ShortLink) error {
	return a.shortLinks.Add(sl)
}

// GetLink returns a short link which matches the given code
func (a *API) GetLink(code string) (shortLink shortlinks.ShortLink, err error) {
	return a.shortLinks.Get(code)
}
