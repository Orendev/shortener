package main

import (
	"github.com/Orendev/shortener/internal/configs"
	"github.com/Orendev/shortener/internal/http"
	storage2 "github.com/Orendev/shortener/internal/storage"
	"log"
)

func main() {
	cfg, err := configs.New()
	if err != nil {
		log.Fatal(err)
		return
	}

	fileDB, err := storage2.NewFileDB(cfg)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = fileDB.Load()
	if err != nil {
		log.Fatal(err)
		return
	}

	memoryStorage, err := storage2.NewMemoryStorage(cfg, fileDB)
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
