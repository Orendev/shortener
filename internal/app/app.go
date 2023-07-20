package app

import (
	"context"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/Orendev/shortener/internal/config"
	"github.com/Orendev/shortener/internal/logger"
	"github.com/Orendev/shortener/internal/repository"
	"github.com/Orendev/shortener/internal/repository/memory"
	"github.com/Orendev/shortener/internal/repository/postgres"
	"github.com/Orendev/shortener/internal/routes"
)

type App struct {
	repo repository.Storage
}

var shutdownTimeout = 10 * time.Second

func Run(cfg *config.Configs) {
	ctx := context.Background()

	var repo repository.Storage

	if len(cfg.Database.DatabaseDSN) > 0 {
		pg, err := postgres.NewRepository(cfg.Database.DatabaseDSN)
		if err != nil {
			log.Fatal(err)
			return
		}

		shutdownCtx, cancel := context.WithTimeout(ctx, shutdownTimeout)
		defer cancel()
		err = pg.Bootstrap(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}

		repo = pg

	} else {

		mem, err := memory.NewRepository(cfg.File.FileStoragePath)
		if err != nil {
			logger.Log.Sugar().Errorf("error memory init: %s", err)
		}

		repo = mem
	}

	defer func() {
		err := repo.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	a := NewApp(repo)

	a.startServer(&http.Server{
		Addr:    cfg.Server.Addr,
		Handler: routes.Router(a.repo, cfg.BaseURL),
	})
}

func NewApp(repo repository.Storage) *App {
	return &App{repo: repo}
}

func (a *App) startServer(srv *http.Server) {
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalf("failed to start server %s", err)
	}
}
