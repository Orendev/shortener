package handlers

import (
	"encoding/json"
	models "github.com/Orendev/shortener/internal/app/models/shortlink"
	repository "github.com/Orendev/shortener/internal/app/storage"
	"github.com/Orendev/shortener/internal/random"
	"io"
	"net/http"
	"strings"
)

type Handler struct {
	ShortLinkRepository repository.ShortLinkRepository
}

func NewHandler(repository repository.ShortLinkRepository) Handler {
	return Handler{ShortLinkRepository: repository}
}

func (h *Handler) ShortLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	code := strings.TrimPrefix(r.URL.Path, "/")

	if shortLink, err := h.ShortLinkRepository.Get(code); err == nil {
		w.Header().Add("Location", shortLink.Link)
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

	sl := models.ShortLink{
		Code: random.Strn(8),
		Link: string(body),
	}

	if _, err := h.ShortLinkRepository.Add(&sl); err != nil {
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
	w.Header().Set("Content-Type", "application/json")
	var shortLink models.ShortLink
	var req models.ShortLinkRequest
	dec := json.NewDecoder(r.Body)
	// читаем тело запроса и декодируем
	if err := dec.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Генерируем уникальный код
	shortLink.Code = random.Strn(8)
	// Сохраняем url
	shortLink.Link = req.URL

	// Сохраним модель
	_, err := h.ShortLinkRepository.Add(&shortLink)

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
