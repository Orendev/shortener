package configs

import (
	"flag"
	models "github.com/Orendev/shortener/internal/app/models/shortlink"
	"github.com/caarlos0/env/v8"
	"log"
)

var cfg Configs = Configs{}
var addr string
var baseURL string
var flagLogLevel string

type Configs struct {
	Addr         string `env:"SERVER_ADDRESS"`
	Host         string `env:"HOST"`
	Port         string `env:"PORT"`
	BaseURL      string `env:"BASE_URL"`
	Memory       map[string]models.ShortLink
	FlagLogLevel string `env:"FLAG_LOG_LEVEL"`
}

func New() (*Configs, error) {

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	flag.StringVar(&addr, "a", "localhost:8080", "Адрес запуска сервера localhost:8080")
	flag.StringVar(&baseURL, "b", "http://localhost:8080", "Базовый URL http://localhost:8080")
	flag.StringVar(&flagLogLevel, "ll", "info", "log level")
	flag.Parse()

	if len(cfg.Addr) == 0 {
		cfg.Addr = addr
	}
	if len(cfg.BaseURL) == 0 {
		cfg.BaseURL = baseURL
	}

	if len(cfg.FlagLogLevel) == 0 {
		cfg.FlagLogLevel = flagLogLevel
	}

	cfg.Memory = map[string]models.ShortLink{}

	return &cfg, nil
}
