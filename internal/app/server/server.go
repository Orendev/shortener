package server

import (
	"fmt"
	"github.com/Orendev/shortener/internal/api"
	"net/http"
)

type Server struct {
	server *http.Server
	cfg    *Config
}

type Config struct {
	Host string
	Port string
}

func New(cfg *Config, a *api.API) (*Server, error) {

	h := &Handlers{api: a}
	r := h.routes()

	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
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
