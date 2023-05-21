package models

import "errors"

// ShortLink модель коротких ссылок
type ShortLink struct {
	UUID        string `json:"uuid"`
	UserID      string `json:"user_id"`
	Code        string `json:"code"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
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

func (sl ShortLinkRequest) Validate() error {
	var err error
	if sl.URL == "" {
		err = errors.New("the URL field is required")
	}

	return err
}
