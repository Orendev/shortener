package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Orendev/shortener/internal/models"
	"github.com/Orendev/shortener/internal/random"
	"github.com/Orendev/shortener/internal/storage"
	"github.com/google/uuid"
	"io"
	"net/http"
	"strings"
	"time"
)

type Handler struct {
	shortLinkStorage storage.ShortLinkStorage
	baseURL          string
}

func NewHandler(storage storage.ShortLinkStorage, baseUrl string) Handler {
	return Handler{shortLinkStorage: storage, baseURL: baseUrl}
}

func (h *Handler) ShortLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	code := strings.TrimPrefix(r.URL.Path, "/")

	if shortLink, err := h.shortLinkStorage.GetByCode(r.Context(), code); err == nil {
		w.Header().Add("Location", shortLink.OriginalURL)
		w.WriteHeader(http.StatusTemporaryRedirect)

		return
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (h *Handler) ShortLinkAdd(w http.ResponseWriter, r *http.Request) {

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
	code := random.Strn(8)
	shortLink := models.ShortLink{
		UUID:        uuid.New().String(),
		Code:        code,
		OriginalURL: req.URL,
		ShortURL:    fmt.Sprintf("%s/%s", strings.TrimPrefix(h.baseURL, "/"), code),
	}

	if _, err = h.shortLinkStorage.Add(r.Context(), &shortLink); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)

	_, err = w.Write([]byte(shortLink.ShortURL))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}

func (h *Handler) APIShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req models.ShortLinkRequest
	dec := json.NewDecoder(r.Body)
	// читаем тело запроса и декодируем
	if err := dec.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	code := random.Strn(8)
	shortLink := models.ShortLink{
		UUID:        uuid.New().String(),
		Code:        code,
		OriginalURL: req.URL,
		ShortURL:    fmt.Sprintf("%s/%s", strings.TrimPrefix(h.baseURL, "/"), code),
	}

	// Сохраним модель
	_, err := h.shortLinkStorage.Add(r.Context(), &shortLink)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// заполняем модель ответа
	resp := models.ShortLinkResponse{
		Result: shortLink.ShortURL,
	}

	enc, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	_, err = w.Write(enc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (h Handler) Ping(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	err := h.shortLinkStorage.Ping(ctx)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
