package routes

import (
	service "github.com/Orendev/shortener/internal/app"
	"github.com/Orendev/shortener/internal/configs"
	"github.com/Orendev/shortener/internal/handlers"
	"github.com/Orendev/shortener/internal/logger"
	"github.com/Orendev/shortener/internal/middlewares"
	"github.com/Orendev/shortener/internal/storage"
	"github.com/go-chi/chi/v5"
)

func Routes(router chi.Router, storage storage.ShortLinkStorage, cfg *configs.Configs) chi.Router {

	h := handlers.NewHandler(service.NewService(storage, cfg))

	if err := logger.NewLogger(cfg.FlagLogLevel); err != nil {
		panic(err)
	}

	router.Use(middlewares.Logger)
	router.Use(middlewares.Gzip)

	router.Route("/", func(r chi.Router) {
		r.Get("/{id}", h.ShortLink)
		r.Post("/", h.ShortLinkAdd)
		r.Post("/api/shorten", h.APIShorten)
	})

	return router
}
