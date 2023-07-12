package routes

import (
	"github.com/Orendev/shortener/internal/handler"
	"github.com/Orendev/shortener/internal/middlewares"
	"github.com/Orendev/shortener/internal/repository"
	"github.com/go-chi/chi/v5"
)

func Router(router chi.Router, repo repository.Storage, baseURL string) chi.Router {

	h := handler.NewHandler(repo, baseURL)

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
