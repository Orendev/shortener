package routes

import (
	"github.com/Orendev/shortener/internal/config"
	"github.com/Orendev/shortener/internal/logger"
	"github.com/Orendev/shortener/internal/middlewares"
	"github.com/Orendev/shortener/internal/storage"
	"github.com/Orendev/shortener/internal/transport/http"
	"github.com/go-chi/chi/v5"
)

func Routes(router chi.Router, storage storage.ShortLinkStorage, cfg *config.Configs) chi.Router {

	h := http.NewHandler(storage, cfg.BaseURL)

	if err := logger.NewLogger(cfg.Log.FlagLogLevel); err != nil {
		panic(err)
	}

	router.Use(middlewares.Logger)
	router.Use(middlewares.Gzip)

	router.Route("/", func(r chi.Router) {
		r.Get("/{id}", h.ShortLink)
		r.Get("/ping", h.Ping)
		r.Post("/", h.ShortLinkAdd)
		r.Post("/api/shorten", h.Shorten)
		r.Post("/api/shorten/batch", h.ShortenBatch)
	})

	return router
}
