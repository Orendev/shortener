package configs

import (
	"flag"
	"github.com/Orendev/shortener/internal/app/repository/shortlinks/Model"
	"github.com/caarlos0/env/v8"
	"log"
)

var cfg Configs = Configs{}
var addr string
var baseURL string

type Configs struct {
	Addr    string `env:"SERVER_ADDRESS"`
	Host    string `env:"HOST"`
	Port    string `env:"PORT"`
	BaseURL string `env:"BASE_URL"`
	Memory  map[string]Model.ShortLink
}

func New() (*Configs, error) {

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	flag.StringVar(&addr, "a", "localhost:8080", "Адрес запуска сервера localhost:8080")
	flag.StringVar(&baseURL, "b", "http://localhost:8080", "Базовый URL http://localhost:8080")
	flag.Parse()

	if len(cfg.Addr) == 0 {
		cfg.Addr = addr
	}
	if len(cfg.BaseURL) == 0 {
		cfg.BaseURL = baseURL
	}

	cfg.Memory = map[string]Model.ShortLink{}

	return &cfg, nil
}
