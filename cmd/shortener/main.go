package main

import (
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

	postgresStorage, err := storage.NewPostgresStorage(cfg.DatabaseDSN)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer func() {
		err := postgresStorage.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	memoryStorage, err := storage.NewMemoryStorage(cfg, postgresStorage, file)
	if err != nil {
		log.Fatal(err)
		return
	}

	srv, err := app.NewServer(cfg, services.NewService(memoryStorage, cfg))
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Fatal(srv.Start())

}
