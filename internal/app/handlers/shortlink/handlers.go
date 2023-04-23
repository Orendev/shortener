package shortlink

import (
	"encoding/json"
	models "github.com/Orendev/shortener/internal/app/models/shortlink"
	repository "github.com/Orendev/shortener/internal/app/repository/shortlink"
	service "github.com/Orendev/shortener/internal/app/service/shortlink"
	"github.com/Orendev/shortener/internal/configs"
	"github.com/Orendev/shortener/internal/random"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"strings"
)

type handler struct {
	ShortLinkRepository repository.ShortLinkRepository
}

func newHandler(repository repository.ShortLinkRepository) handler {
	return handler{ShortLinkRepository: repository}
}

func Routes(router chi.Router, cfg *configs.Configs) chi.Router {

	memoryStorage, err := repository.NewMemoryStorage(cfg)
	if err != nil {
		return nil
	}

	h := newHandler(service.NewService(memoryStorage, cfg))

	router.Route("/", func(r chi.Router) {
		r.Get("/{id}", h.handleShortLink)
		r.Post("/", h.handleShortLinkAdd)
		r.Post("/api/shorten", h.handleApiShorten)
	})

	return router
}

func (h *handler) handleShortLink(w http.ResponseWriter, r *http.Request) {
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

func (h *handler) handleShortLinkAdd(w http.ResponseWriter, r *http.Request) {

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

func (h handler) handleApiShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var shortLink models.ShortLink
	var req models.Request
	dec := json.NewDecoder(r.Body)
	// читаем тело запроса и декодируем
	if err := dec.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Генерируем уникальный код
	shortLink.Code = random.Strn(8)
	// Сохраняем url
	shortLink.Link = req.Url

	// Сохраним модель
	_, err := h.ShortLinkRepository.Add(&shortLink)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// заполняем модель ответа
	resp := models.Response{
		Result: shortLink.Result,
	}

	enc, err := json.MarshalIndent(resp, "", "   ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(enc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
