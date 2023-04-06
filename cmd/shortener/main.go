package main

import (
	"github.com/Orendev/shortener/internal/api"
	"github.com/Orendev/shortener/internal/app/repository/shortlinks"
	"github.com/Orendev/shortener/internal/app/server"
	"github.com/Orendev/shortener/internal/configs"
	"log"
)

func main() {

	data := map[string]shortlinks.ShortLink{}
	shortLinkStore, err := shortlinks.NewStorage(data)
	if err != nil {
		return
	}

	sl, err := shortlinks.New(*shortLinkStore)
	if err != nil {
		return
	}

	cfg, err := configs.New()
	if err != nil {
		return
	}

	a, err := api.New(cfg, sl)
	if err != nil {
		return
	}

	srv, err := server.New(cfg, a)
	if err != nil {
		return
	}

	log.Fatal(srv.Start())

}
