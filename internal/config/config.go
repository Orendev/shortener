package config

import (
	"flag"

	"github.com/caarlos0/env/v8"
)

var cfg Configs = Configs{}
var addr string
var baseURL string
var flagLogLevel string
var fileStoragePath string
var databaseDSN string

type Server struct {
	Addr string `env:"SERVER_ADDRESS"`
	Host string `env:"HOST"`
	Port string `env:"PORT"`
}

type File struct {
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

type Log struct {
	FlagLogLevel string `env:"FLAG_LOG_LEVEL"`
}

type Database struct {
	DatabaseDSN string `env:"DATABASE_DSN"`
}

type Configs struct {
	Database Database
	Server   struct {
		Addr string `env:"SERVER_ADDRESS"`
		Host string `env:"HOST"`
		Port string `env:"PORT"`
	}
	File    File
	Log     Log
	BaseURL string `env:"BASE_URL"`
}

func New() (*Configs, error) {

	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}
	flag.StringVar(&addr, "a", "localhost:8080", "Адрес запуска сервера localhost:8080")
	flag.StringVar(&baseURL, "b", "http://localhost:8080", "Базовый URL http://localhost:8080")
	flag.StringVar(&flagLogLevel, "ll", "info", "log level")
	flag.StringVar(&fileStoragePath, "f", "/tmp/short-url-db.json", "Полное имя файла")
	//host=localhost user=shortener password=secret dbname=shortener sslmode=disable
	flag.StringVar(&databaseDSN, "d", "", "Строка с адресом подключения")
	flag.Parse()

	if len(cfg.Server.Addr) == 0 {
		cfg.Server.Addr = addr
	}
	if len(cfg.BaseURL) == 0 {
		cfg.BaseURL = baseURL
	}

	if len(cfg.Log.FlagLogLevel) == 0 {
		cfg.Log.FlagLogLevel = flagLogLevel
	}

	if len(cfg.File.FileStoragePath) == 0 {
		cfg.File.FileStoragePath = fileStoragePath
	}

	if len(cfg.Database.DatabaseDSN) == 0 {
		cfg.Database.DatabaseDSN = databaseDSN
	}

	return &cfg, nil
}
