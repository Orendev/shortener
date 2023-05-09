package main

import (
	service "github.com/Orendev/shortener/internal/app"
	"github.com/Orendev/shortener/internal/configs"
	"github.com/Orendev/shortener/internal/http"
	"github.com/Orendev/shortener/internal/storage"
	"log"
)

func main() {
	cfg, err := configs.New()
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

	srv, err := http.New(cfg, service.NewService(memoryStorage, cfg))
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Fatal(srv.Start())

}
