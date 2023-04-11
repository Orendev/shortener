package http

import (
	"fmt"
	shortLinksHandlers "github.com/Orendev/shortener/internal/app/repository/shortlink/handlers"
	"github.com/Orendev/shortener/internal/configs"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Server struct {
	server *http.Server
	cfg    *configs.Configs
}

func New(cfg *configs.Configs) (*Server, error) {

	r := shortLinksHandlers.Routes(chi.NewRouter(), cfg)

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
