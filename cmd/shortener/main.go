package main

import (
	"github.com/Orendev/shortener/internal/app/http"
	"github.com/Orendev/shortener/internal/configs"
	"log"
)

func main() {
	cfg, err := configs.New()
	if err != nil {
		return
	}

	srv, err := http.New(cfg)
	if err != nil {
		return
	}

	log.Fatal(srv.Start())

}
