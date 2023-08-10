package config

import (
	"flag"

	"github.com/caarlos0/env/v8"
)

var cfg Configs = Configs{}
var addr string
var baseURL string
var isHTTPS bool
var flagLogLevel string
var fileStoragePath string
var databaseDSN string
var keyFile string
var certFile string

// Server configuration
type Server struct {
	Addr    string `env:"SERVER_ADDRESS"`
	Host    string `env:"HOST"`
	Port    string `env:"PORT"`
	IsHTTPS bool   `env:"ENABLE_HTTPS"`
}

// File configuration
type File struct {
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

// Log configuration
type Log struct {
	FlagLogLevel string `env:"FLAG_LOG_LEVEL"`
}

// Database configuration
type Database struct {
	DatabaseDSN string `env:"DATABASE_DSN"`
}

// Cert configuration
type Cert struct {
	CertFile string `env:"FILE_CERT"`
	KeyFile  string `env:"FILE_PRIVATE_KEY"`
}

// Configs configuration
type Configs struct {
	Database Database
	Server   Server
	Cert     Cert
	File     File
	Log      Log
	BaseURL  string `env:"BASE_URL"`
}

// New constructor a new instance of Configs
func New() (*Configs, error) {

	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}
	flag.StringVar(&addr, "a", "localhost:8080", "Адрес запуска сервера localhost:8080")
	flag.StringVar(&baseURL, "b", "http://localhost:8080", "Базовый URL http://localhost:8080")
	flag.BoolVar(&isHTTPS, "s", false, "Включения HTTPS в веб-сервере.")
	flag.StringVar(&flagLogLevel, "ll", "info", "log level")
	flag.StringVar(&fileStoragePath, "f", "/tmp/short-url-db.json", "Полное имя файла")
	flag.StringVar(&keyFile, "fc", "key.pem", "Закрытый ключ")
	flag.StringVar(&certFile, "fk", "cert.pem", "Подписанный центром сертификации, файл сертификата")
	//host=localhost user=shortener password=secret dbname=shortener sslmode=disable
	flag.StringVar(&databaseDSN, "d", "", "Строка с адресом подключения")
	flag.Parse()

	if len(cfg.Server.Addr) == 0 {
		cfg.Server.Addr = addr
	}
	if len(cfg.BaseURL) == 0 {
		cfg.BaseURL = baseURL
	}

	if !cfg.Server.IsHTTPS {
		cfg.Server.IsHTTPS = isHTTPS
	}

	if len(cfg.Cert.CertFile) == 0 {
		cfg.Cert.CertFile = certFile
	}

	if len(cfg.Cert.KeyFile) == 0 {
		cfg.Cert.KeyFile = keyFile
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
