package handlers

import (
	"github.com/Orendev/shortener/internal/repository"
)

// Handler - structure describing the handler
type Handler struct {
	repo    repository.Storage
	baseURL string
}

// NewHandler конструктор создает структуру Handler
func NewHandler(repo repository.Storage, baseURL string) Handler {
	return Handler{repo: repo, baseURL: baseURL}
}
