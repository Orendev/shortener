package config

import (
	"flag"
	"github.com/Orendev/shortener/internal/models"
	"github.com/caarlos0/env/v8"
	"log"
)

var cfg Configs = Configs{}
var addr string
var baseURL string
var flagLogLevel string
var fileStoragePath string
var databaseDSN string

type Configs struct {
	Addr            string `env:"SERVER_ADDRESS"`
	Host            string `env:"HOST"`
	Port            string `env:"PORT"`
	BaseURL         string `env:"BASE_URL"`
	Memory          map[string]models.ShortLink
	FlagLogLevel    string `env:"FLAG_LOG_LEVEL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
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
	flag.StringVar(&fileStoragePath, "f", "/tmp/short-url-db.json", "Полное имя файла")
	//host=localhost user=shortener password=secret dbname=shortener sslmode=disable
	flag.StringVar(&databaseDSN, "d", "", "Строка с адресом подключения")
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

	if len(cfg.FileStoragePath) == 0 {
		cfg.FileStoragePath = fileStoragePath
	}

	if len(cfg.DatabaseDSN) == 0 {
		cfg.DatabaseDSN = databaseDSN
	}

	cfg.Memory = map[string]models.ShortLink{}

	return &cfg, nil
}
