package main

import (
	"github.com/Orendev/shortener/internal/app/http"
	"github.com/Orendev/shortener/internal/app/storage"
	"github.com/Orendev/shortener/internal/configs"
	"log"
)

func main() {
	cfg, err := configs.New()
	if err != nil {
		log.Fatal(err)
		return
	}

	fileDB, err := storage.NewFileDB(cfg)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = fileDB.Load()
	if err != nil {
		log.Fatal(err)
		return
	}

	memoryStorage, err := storage.NewMemoryStorage(cfg, fileDB)
	if err != nil {
		log.Fatal(err)
		return
	}

	srv, err := http.New(cfg, memoryStorage)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Fatal(srv.Start())

}
