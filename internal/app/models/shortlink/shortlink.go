package shortlink

// ShortLink модель коротких ссылок
type ShortLink struct {
	Code   string `json:"-"`
	Link   string `json:"url"`
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
	UUID        uint   `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
