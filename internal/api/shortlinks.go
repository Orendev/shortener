package api

import (
	"fmt"
	"github.com/Orendev/shortener/internal/app/repository/shortlinks"
	"strings"
)

// AddLink add a new short link
func (a *API) AddLink(sl shortlinks.ShortLink) (url string, err error) {
	url = fmt.Sprintf("%s/%s", strings.TrimPrefix(a.cfg.BaseURL, "/"), sl.Code)
	err = a.shortLinks.Add(sl)
	return
}

// GetLink returns a short link which matches the given code
func (a *API) GetLink(code string) (shortLink shortlinks.ShortLink, err error) {
	return a.shortLinks.Get(code)
}
