package models

import "errors"

// ShortLink the short link model.
type ShortLink struct {
	UUID        string `json:"uuid" db:"id"`
	UserID      string `json:"user_id" db:"user_id"`
	Code        string `json:"code" db:"-"`
	ShortURL    string `json:"short_url" db:"short_url"`
	OriginalURL string `json:"original_url" db:"original_url"`
	DeletedFlag bool   `json:"is_deleted" db:"is_deleted"`
}

// ShortLinkResponse describes the server response.
type ShortLinkResponse struct {
	Result string `json:"result"`
}

// ShortLinkRequest describes the client's request.
type ShortLinkRequest struct {
	URL string `json:"url"`
}

// ShortLinkBatchRequest describes the client's request.
type ShortLinkBatchRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// ShortLinkBatchResponse describes the response of the short link list server.
type ShortLinkBatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// ShortLinkUserResponse describes the response of the user's short link server.
type ShortLinkUserResponse struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}

// Validate validation of the input request.
func (sl ShortLinkRequest) Validate() error {
	var err error
	if sl.URL == "" {
		err = errors.New("the URL field is required")
	}

	return err
}
