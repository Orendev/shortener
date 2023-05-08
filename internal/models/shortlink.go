package models

import "errors"

// ShortLink модель коротких ссылок
type ShortLink struct {
	Code   string `json:"-"`
	URL    string `json:"url"`
	Result string `json:"result"`
}

// ShortLinkResponse описывает ответ сервера.
type ShortLinkResponse struct {
	Result string `json:"result"`
}

// ShortLinkRequest описывает запрос клиента.
type ShortLinkRequest struct {
	URL string `json:"url"`
}

type FileDB struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func (sl ShortLinkRequest) Validate() error {
	var err error
	if sl.URL == "" {
		err = errors.New("the URL field is required")
	}

	return err
}
