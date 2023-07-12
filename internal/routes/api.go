package routes

import (
	"log"

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
		log.Fatal(err)
	}

	router.Use(middlewares.Logger)
	router.Use(middlewares.Gzip)
	router.Use(middlewares.Auth)

	router.Route("/", func(r chi.Router) {
		r.Get("/{id}", h.ShortLink)
		r.Get("/ping", h.Ping)
		r.Get("/api/user/urls", h.UserUrls)
		r.Post("/", h.ShortLinkAdd)
		r.Post("/api/shorten", h.Shorten)
		r.Post("/api/shorten/batch", h.ShortenBatch)
		r.Delete("/api/user/urls", h.DeleteUserUrls)
	})

	return router
}
