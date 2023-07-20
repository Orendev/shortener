package main

import (
	"log"
	_ "net/http/pprof"

	"github.com/Orendev/shortener/internal/app"
	"github.com/Orendev/shortener/internal/config"
	"github.com/Orendev/shortener/internal/logger"
)

func main() {

	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	if err := logger.NewLogger("info"); err != nil {
		log.Fatal(err)
	}

	app.Run(cfg)
}
