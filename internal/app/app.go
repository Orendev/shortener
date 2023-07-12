package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Orendev/shortener/internal/config"
	"github.com/Orendev/shortener/internal/repository"
	"github.com/Orendev/shortener/internal/repository/memory"
	"github.com/Orendev/shortener/internal/repository/postgres"
	"github.com/Orendev/shortener/internal/routes"
	"github.com/go-chi/chi/v5"
)

type App struct {
	repo repository.Storage
}

var shutdownTimeout = 10 * time.Second

func Run(cfg *config.Configs) {
	ctx := gracefulShutdown()

	var repo repository.Storage
	if len(cfg.Database.DatabaseDSN) > 0 {
		pg, err := postgres.NewRepository(cfg.Database.DatabaseDSN)
		if err != nil {
			log.Fatal(err)
			return
		}

		defer func() {
			err = pg.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		err = pg.Bootstrap(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}

		repo = pg

	} else {
		file, err := repository.NewFile(cfg)
		if err != nil {
			log.Fatal(err)
			return
		}

		repo, err = memory.NewRepository(cfg.Memory, file)
		if err != nil {
			log.Fatal(err)
			return
		}
	}

	defer func() {
		err := repo.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	a := NewApp(ctx, repo)

	a.startServer(ctx, &http.Server{
		Addr:    cfg.Server.Addr,
		Handler: routes.Router(chi.NewRouter(), a.repo, cfg.BaseURL),
	})
}

func NewApp(_ context.Context, repo repository.Storage) *App {
	return &App{repo: repo}
}

func (a *App) startServer(ctx context.Context, srv *http.Server) {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := srv.ListenAndServe()
		if err != nil {
			log.Fatalf("failed to start server %s", err)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err := srv.Shutdown(shutdownCtx)
	if err != nil {
		log.Fatalf("failed to shudown server %s", err)
	}

	wg.Wait()
}

func gracefulShutdown() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	irqSig := make(chan os.Signal, 1)
	// Получено сообщение о завершении работы от операционной системы.
	signal.Notify(irqSig, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		<-irqSig
		cancel()
	}()
	return ctx
}
