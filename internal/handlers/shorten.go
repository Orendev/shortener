package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Orendev/shortener/internal/auth"
	"github.com/Orendev/shortener/internal/models"
	"github.com/Orendev/shortener/internal/random"
	"github.com/Orendev/shortener/internal/repository"
	"github.com/google/uuid"
)

// GetShorten we will get a short link by its code.
func (h *Handler) GetShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	code := strings.TrimPrefix(r.URL.Path, "/")

	shortLink, err := h.repo.GetByCode(r.Context(), code)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if shortLink.DeletedFlag {
		// Целевой запрос больше не доступен
		w.WriteHeader(http.StatusGone)
		return
	}

	w.Header().Add("Location", shortLink.OriginalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// PostShorten save the short link.
func (h *Handler) PostShorten(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/plain")

	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req := models.ShortLinkRequest{}
	req.URL = string(body)

	if err = req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := auth.GetAuthIdentifier(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	code := random.Strn(8)
	shortLink := &models.ShortLink{
		UUID:        uuid.New().String(),
		UserID:      userID,
		Code:        code,
		OriginalURL: req.URL,
		ShortURL:    fmt.Sprintf("%s/%s", strings.TrimPrefix(h.baseURL, "/"), code),
		DeletedFlag: false,
	}

	err = h.repo.Save(r.Context(), *shortLink)

	if err != nil && !errors.Is(err, repository.ErrConflict) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if errors.Is(err, repository.ErrConflict) {
		w.WriteHeader(http.StatusConflict)
		shortLink, err = h.repo.GetByOriginalURL(r.Context(), req.URL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)

	_, err = w.Write([]byte(shortLink.ShortURL))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}
