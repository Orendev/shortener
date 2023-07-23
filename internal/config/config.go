package config

import (
	"flag"

	"github.com/caarlos0/env/v8"
)

var cfg Configs = Configs{}

// Server конфигурация сервера
type Server struct {
	Addr string `env:"SERVER_ADDRESS"`
	Host string `env:"HOST"`
	Port string `env:"PORT"`
}

// File конфигурация файла
type File struct {
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

// Log конфигурация LOG
type Log struct {
	FlagLogLevel string `env:"FLAG_LOG_LEVEL"`
}

// Database конфигурация БД
type Database struct {
	DatabaseDSN string `env:"DATABASE_DSN"`
}

// Configs конфигурация конфигаа
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

// New констуктор
func New() (*Configs, error) {

	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}

	if len(cfg.Server.Addr) == 0 {
		flag.StringVar(&cfg.Server.Addr, "a", "localhost:8080", "Адрес запуска сервера localhost:8080")
	}
	if len(cfg.BaseURL) == 0 {
		flag.StringVar(&cfg.BaseURL, "b", "http://localhost:8080", "Базовый URL http://localhost:8080")
	}

	if len(cfg.Log.FlagLogLevel) == 0 {
		flag.StringVar(&cfg.Log.FlagLogLevel, "ll", "info", "log level")
	}

	if len(cfg.File.FileStoragePath) == 0 {
		flag.StringVar(&cfg.File.FileStoragePath, "f", "/tmp/short-url-db.json", "Полное имя файла")
	}

	if len(cfg.Database.DatabaseDSN) == 0 {
		flag.StringVar(&cfg.Database.DatabaseDSN, "d", "", "Строка с адресом подключения")
	}

	flag.Parse()

	return &cfg, nil
}
