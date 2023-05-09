package app

import (
	"fmt"
	"github.com/Orendev/shortener/internal/config"
	"github.com/Orendev/shortener/internal/routes"
	"github.com/Orendev/shortener/internal/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Server struct {
	server *http.Server
	cfg    *config.Configs
}

func NewServer(cfg *config.Configs, storage storage.ShortLinkStorage) (*Server, error) {

	r := routes.Routes(chi.NewRouter(), storage, cfg)

	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	if len(cfg.Addr) > 0 {
		addr = cfg.Addr
	}

	httpServer := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	return &Server{
		server: httpServer,
		cfg:    cfg,
	}, nil
}

func (srv *Server) Start() error {
	return srv.server.ListenAndServe()
}
