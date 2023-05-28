package models

import "errors"

// ShortLink модель коротких ссылок
type ShortLink struct {
	UUID        string `json:"uuid" db:"id"`
	UserID      string `json:"user_id" db:"user_id"`
	Code        string `json:"code" db:"-"`
	ShortURL    string `json:"short_url" db:"short_url"`
	OriginalURL string `json:"original_url" db:"original_url"`
	DeletedFlag bool   `json:"-" db:"is_deleted"`
}

// ShortLinkResponse описывает ответ сервера.
type ShortLinkResponse struct {
	Result string `json:"result"`
}

// ShortLinkRequest описывает запрос клиента.
type ShortLinkRequest struct {
	URL string `json:"url"`
}

// ShortLinkBatchRequest описывает запрос клиента.
type ShortLinkBatchRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// ShortLinkBatchResponse описывает ответ сервера.
type ShortLinkBatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// ShortLinkUserResponse описывает ответ сервера.
type ShortLinkUserResponse struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}

func (sl ShortLinkRequest) Validate() error {
	var err error
	if sl.URL == "" {
		err = errors.New("the URL field is required")
	}

	return err
}
