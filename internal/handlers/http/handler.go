package http

import (
	"github.com/Orendev/shortener/internal/models"
	"github.com/Orendev/shortener/internal/repository"
)

// Handler - structure describing the handler
type Handler struct {
	repo                  repository.Storage
	baseURL               string
	trustedSubnet         string
	msgDeleteUserUrlsChan chan models.Message
}

// NewHandler конструктор создает структуру Handler
func NewHandler(repo repository.Storage, baseURL, trustedSubnet string) Handler {
	instance := Handler{repo: repo, baseURL: baseURL, msgDeleteUserUrlsChan: make(chan models.Message, 10), trustedSubnet: trustedSubnet}
	// запустим горутину с фоновым удалением пользовательских ссылок
	go instance.flushDeleteShortLink()
	return instance
}
