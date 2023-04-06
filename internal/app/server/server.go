package server

import (
	"fmt"
	"github.com/Orendev/shortener/internal/api"
	"github.com/Orendev/shortener/internal/configs"
	"net/http"
)

type Server struct {
	server *http.Server
	cfg    *configs.Configs
}

func New(cfg *configs.Configs, a *api.API) (*Server, error) {

	h := &Handlers{api: a}
	r := h.routes()

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
