package routes

import (
	"github.com/Orendev/shortener/internal/handler"
	"github.com/Orendev/shortener/internal/middlewares"
	"github.com/Orendev/shortener/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Router(repo repository.Storage, baseURL string) *chi.Mux {

	h := handler.NewHandler(repo, baseURL)
	router := chi.NewRouter()
	router.Use(middlewares.Logger)
	router.Use(middlewares.Gzip)
	router.Use(middlewares.Auth)

	router.Mount("/debug", middleware.Profiler())

	router.Route("/api", func(r chi.Router) {
		r.Get("/user/urls", h.UserUrls)
		r.Post("/shorten", h.Shorten)
		r.Post("/shorten/batch", h.ShortenBatch)
		r.Delete("/user/urls", h.DeleteUserUrls)
	})

	router.Route("/", func(r chi.Router) {
		r.Get("/{id}", h.ShortLink)
		r.Get("/ping", h.Ping)
		r.Post("/", h.ShortLinkAdd)
	})

	return router
}
