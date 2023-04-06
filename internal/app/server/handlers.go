package server

import (
	"github.com/Orendev/shortener/internal/api"
	"github.com/Orendev/shortener/internal/app/repository/shortlinks"
	"github.com/Orendev/shortener/internal/pkg/random"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"strings"
)

type Handlers struct {
	api *api.API
}

func (h *Handlers) routes() chi.Router {
	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		r.Get("/{id}", h.handleShortLink)
		r.Post("/", h.handleShortLinkAdd)
	})
	return router
}

func (h *Handlers) handleShortLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	code := strings.TrimPrefix(r.URL.Path, "/")

	if shortLink, err := h.api.GetLink(code); err == nil {
		w.Header().Add("Location", shortLink.Link)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (h *Handlers) handleShortLinkAdd(w http.ResponseWriter, r *http.Request) {

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

	sl := shortlinks.ShortLink{
		Code: random.Strn(8),
		Link: string(body),
	}

	url, err := h.api.AddLink(sl)

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
