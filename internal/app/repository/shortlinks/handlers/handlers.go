package handlers

import (
	"github.com/Orendev/shortener/internal/app/repository/shortlinks"
	"github.com/Orendev/shortener/internal/app/repository/shortlinks/model"
	"github.com/Orendev/shortener/internal/app/repository/shortlinks/storage"
	"github.com/Orendev/shortener/internal/configs"
	"github.com/Orendev/shortener/internal/pkg/random"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"strings"
)

type handler struct {
	ShortLinkService shortlinks.Service
}

func newHandler(sls shortlinks.Service) handler {
	return handler{ShortLinkService: sls}
}

func Routes(router chi.Router, cfg *configs.Configs) chi.Router {
	s, err := storage.New(cfg)
	if err != nil {
		return nil
	}
	service, err := shortlinks.New(s, cfg)
	if err != nil {
		return nil
	}

	h := newHandler(service)

	router.Route("/", func(r chi.Router) {
		r.Get("/{id}", h.handleShortLink)
		r.Post("/", h.handleShortLinkAdd)
	})

	return router
}

func (h *handler) handleShortLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	code := strings.TrimPrefix(r.URL.Path, "/")

	if shortLink, err := h.ShortLinkService.Get(code); err == nil {
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

	sl := model.ShortLink{
		Code: random.Strn(8),
		Link: string(body),
	}

	url, err := h.ShortLinkService.Add(sl)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)

	_, err = w.Write([]byte(url))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}
