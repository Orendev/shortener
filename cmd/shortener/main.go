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
	var ctx = context.Background()

	dbStore, err := storage.NewPostgresStorage(cfg.DatabaseDSN)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer func() {
		err := dbStore.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	if len(cfg.DatabaseDSN) > 0 {
		store = dbStore
		err = dbStore.Bootstrap(ctx)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		store, err = storage.NewMemoryStorage(cfg, dbStore, file)
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
