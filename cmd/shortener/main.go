package main

import (
	"context"
	"log"
	"time"

	"github.com/Orendev/shortener/internal/app"
	"github.com/Orendev/shortener/internal/config"
	"github.com/Orendev/shortener/internal/services"
	"github.com/Orendev/shortener/internal/storage"
)

const (
	shutdownTimeout = 5 * time.Second
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
		return
	}

	var store storage.ShortLinkStorage

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if len(cfg.Database.DatabaseDSN) > 0 {
		pg, err := storage.NewPostgresStorage(cfg.Database.DatabaseDSN)
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

		err = pg.Bootstrap(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}

		store = pg

	} else {
		file, err := storage.NewFile(cfg)
		if err != nil {
			log.Fatal(err)
			return
		}

		store, err = storage.NewMemoryStorage(cfg.Memory, file)
		if err != nil {
			log.Fatal(err)
			return
		}
	}

	srv, err := app.NewServer(cfg, services.NewService(store))
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Fatal(srv.Start())

}
