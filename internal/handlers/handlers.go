package handlers

import (
	"encoding/json"
	"github.com/Orendev/shortener/internal/models"
	"github.com/Orendev/shortener/internal/random"
	"github.com/Orendev/shortener/internal/storage"
	"io"
	"net/http"
	"strings"
)

type Handler struct {
	shortLinkStorage storage.ShortLinkStorage
}

func NewHandler(storage storage.ShortLinkStorage) Handler {
	return Handler{shortLinkStorage: storage}
}

func (h *Handler) ShortLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	code := strings.TrimPrefix(r.URL.Path, "/")

	if shortLink, err := h.shortLinkStorage.Get(code); err == nil {
		w.Header().Add("Location", shortLink.URL)
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

	sl := models.ShortLink{
		Code: random.Strn(8),
		URL:  req.URL,
	}

	if _, err = h.shortLinkStorage.Add(&sl); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)

	_, err = w.Write([]byte(sl.Result))

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

	var shortLink models.ShortLink
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

	// Генерируем уникальный код
	shortLink.Code = random.Strn(8)
	// Сохраняем url
	shortLink.URL = req.URL

	// Сохраним модель
	_, err := h.shortLinkStorage.Add(&shortLink)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// заполняем модель ответа
	resp := models.ShortLinkResponse{
		Result: shortLink.Result,
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
