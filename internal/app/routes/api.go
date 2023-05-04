package routes

import (
	"github.com/Orendev/shortener/internal/app/handlers"
	"github.com/Orendev/shortener/internal/app/middlewares"
	service "github.com/Orendev/shortener/internal/app/service/shortlink"
	repository "github.com/Orendev/shortener/internal/app/storage"
	"github.com/Orendev/shortener/internal/configs"
	"github.com/Orendev/shortener/internal/logger"
	"github.com/go-chi/chi/v5"
)

func Routes(router chi.Router, repository repository.ShortLinkRepository, cfg *configs.Configs) chi.Router {

	h := handlers.NewHandler(service.NewService(repository, cfg))

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
