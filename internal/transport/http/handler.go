package http

import (
	"context"
	"encoding/json"
	"errors"
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

func NewHandler(storage storage.ShortLinkStorage, baseURL string) Handler {
	return Handler{shortLinkStorage: storage, baseURL: baseURL}
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

	if err = h.shortLinkStorage.Save(r.Context(), shortLink); err != nil {
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

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
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
	shortLink := &models.ShortLink{
		UUID:        uuid.New().String(),
		Code:        code,
		OriginalURL: req.URL,
		ShortURL:    fmt.Sprintf("%s/%s", strings.TrimPrefix(h.baseURL, "/"), code),
	}

	// Сохраним модель
	err := h.shortLinkStorage.Save(r.Context(), *shortLink)

	if err != nil && !errors.Is(err, storage.ErrConflict) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if errors.Is(err, storage.ErrConflict) {
		w.WriteHeader(http.StatusConflict)
		shortLink, err = h.shortLinkStorage.GetByOriginalUrl(r.Context(), req.URL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
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

func (h *Handler) ShortenBatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var reqData []models.ShortLinkBatchRequest
	dec := json.NewDecoder(r.Body)
	// читаем тело запроса и декодируем
	if err := dec.Decode(&reqData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	shortLinksInsert := make([]models.ShortLink, 0, len(reqData))
	shortLinksUpdate := make([]models.ShortLink, 0, len(reqData))
	shortLinkBatchResponse := make([]models.ShortLinkBatchResponse, 0, len(reqData))

	for _, req := range reqData {
		code := random.Strn(8)
		var model *models.ShortLink

		model, err := h.shortLinkStorage.GetByID(r.Context(), req.CorrelationID)

		if err != nil {
			model = &models.ShortLink{
				UUID:        req.CorrelationID,
				Code:        code,
				OriginalURL: req.OriginalURL,
				ShortURL:    fmt.Sprintf("%s/%s", strings.TrimPrefix(h.baseURL, "/"), code),
			}

			shortLinksInsert = append(shortLinksInsert, *model)

		} else {

			model.OriginalURL = req.OriginalURL
			shortLinksUpdate = append(shortLinksUpdate, *model)
		}

		// заполняем модель ответа
		shortLinkBatchResponse = append(shortLinkBatchResponse, models.ShortLinkBatchResponse{
			CorrelationID: model.UUID,
			ShortURL:      model.ShortURL,
		})
	}

	// Сохраним модель
	err := h.shortLinkStorage.InsertBatch(r.Context(), shortLinksInsert)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.shortLinkStorage.UpdateBatch(r.Context(), shortLinksUpdate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// заполняем модель ответа
	enc, err := json.Marshal(shortLinkBatchResponse)
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
