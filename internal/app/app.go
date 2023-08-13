package app

import (
	"context"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Orendev/shortener/internal/config"
	"github.com/Orendev/shortener/internal/logger"
	"github.com/Orendev/shortener/internal/repository"
	"github.com/Orendev/shortener/internal/repository/memory"
	"github.com/Orendev/shortener/internal/repository/postgres"
	"github.com/Orendev/shortener/internal/routes"
	"github.com/Orendev/shortener/internal/tls"
	"go.uber.org/zap"
)

// App - structure describing the application
type App struct {
	repo repository.Storage
}

var shutdownTimeout = 10 * time.Second

// Run starts the application.
func Run(cfg *config.Configs) {
	ctx := gracefulShutdown()

	var repo repository.Storage

	if len(cfg.Database.DatabaseDSN) > 0 {
		pg, err := postgres.NewPostgres(cfg.Database.DatabaseDSN)
		if err != nil {
			logger.Log.Sugar().Errorf("error postgres init: %s", err)
			return
		}

		shutdownCtx, cancel := context.WithTimeout(ctx, shutdownTimeout)
		defer cancel()
		err = pg.Bootstrap(shutdownCtx)
		if err != nil {
			logger.Log.Sugar().Errorf("error postgres shutdownCtx: %s", err)
		}

		repo = pg

	} else {

		mem, err := memory.NewMemory(cfg.File.FileStoragePath)
		if err != nil {
			logger.Log.Sugar().Errorf("error memory init: %s", err)
		}

		repo = mem
	}

	defer func() {
		err := repo.Close()
		if err != nil {
			logger.Log.Error("error repo", zap.Error(err))
		}
	}()

	a := NewApp(repo)

	err := tls.New(cfg.Cert.CertFile, cfg.Cert.KeyFile)
	if err != nil {
		logger.Log.Error("error tls init", zap.Error(err))
	}

	a.startServer(ctx, &http.Server{
		Addr:    cfg.Server.Addr,
		Handler: routes.Router(a.repo, cfg.BaseURL),
	},
		cfg.Server.IsHTTPS,
		cfg.Cert.CertFile,
		cfg.Cert.KeyFile,
	)
}

// NewApp constructor for the application.
func NewApp(repo repository.Storage) *App {
	return &App{repo: repo}
}

func (a *App) startServer(ctx context.Context, srv *http.Server, isHTTPS bool, certFile, keyFile string) {
	var err error
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if isHTTPS {
			err = srv.ListenAndServeTLS(certFile, keyFile)
		} else {
			err = srv.ListenAndServe()
		}
		if err != nil {
			log.Fatalf("failed to start server %s", err)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = srv.Shutdown(shutdownCtx)
	if err != nil {
		log.Fatalf("failed to shudown server %s", err)
	}

	wg.Wait()

}

func gracefulShutdown() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	irqSig := make(chan os.Signal, 1)
	// Получено сообщение о завершении работы от операционной системы.
	signal.Notify(irqSig, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	go func() {
		<-irqSig
		cancel()
	}()
	return ctx
}
