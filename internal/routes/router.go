package routes

import (
	"github.com/Orendev/shortener/internal/handlers/http"
	middlewares "github.com/Orendev/shortener/internal/middlewares/http"
	"github.com/Orendev/shortener/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Router api handlers
func Router(repo repository.Storage, baseURL, trustedSubnet string) *chi.Mux {

	h := http.NewHandler(repo, baseURL, trustedSubnet)
	router := chi.NewRouter()
	router.Use(middlewares.Logger)
	router.Use(middlewares.Gzip)
	router.Use(middlewares.Auth)

	router.Mount("/debug", middleware.Profiler())

	router.Route("/api", func(r chi.Router) {
		r.Get("/user/urls", h.GetAPIUserUrls)
		r.Get("/internal/stats", h.GetAPIStats)
		r.Post("/shorten", h.PostAPIShorten)
		r.Post("/shorten/batch", h.PostAPIShortenBatch)
		r.Delete("/user/urls", h.DeleteAPIUserUrls)
	})

	router.Route("/", func(r chi.Router) {
		r.Get("/{id}", h.GetShorten)
		r.Get("/ping", h.GetPing)
		r.Post("/", h.PostShorten)
	})

	return router
}
