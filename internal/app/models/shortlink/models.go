package shortlink

// ShortLink модель коротких ссылок
type ShortLink struct {
	Code   string `json:"-"`
	Link   string `json:"url"`
	Result string `json:"result"`
}

// Response описывает ответ сервера.
type Response struct {
	Result string `json:"result"`
}

// Request описывает запрос клиента.
type Request struct {
	URL string `json:"url"`
}
