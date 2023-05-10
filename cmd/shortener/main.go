package main

import (
	"context"
	"github.com/Orendev/shortener/internal/app"
	"github.com/Orendev/shortener/internal/config"
	"github.com/Orendev/shortener/internal/services"
	"github.com/Orendev/shortener/internal/storage"
	"log"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
		return
	}

	file, err := storage.NewFile(cfg)
	if err != nil {
		log.Fatal(err)
		return
	}

	var store storage.ShortLinkStorage

	if len(cfg.DatabaseDSN) > 0 {

		store, err = storage.NewPostgresStorage(context.Background(), cfg.DatabaseDSN)
		if err != nil {
			log.Fatal(err)
			return
		}

		defer func() {
			err := store.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()

	} else {
		store, err = storage.NewMemoryStorage(cfg, nil, file)
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
