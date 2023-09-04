package config

import (
	"encoding/json"
	"flag"
	"os"
	"strconv"
)

// Server configuration
type Server struct {
	Addr    string `env:"SERVER_ADDRESS"`
	IsHTTPS bool   `env:"ENABLE_HTTPS"`
}

type GRPCServer struct {
	Addr string `env:"GRPC_ADDRESS"`
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

// Configs configuration application
type Configs struct {
	Database      Database
	Server        Server
	GRPC          GRPCServer
	Cert          Cert
	File          File
	Log           Log
	BaseURL       string `env:"BASE_URL"`
	Config        string `env:"CONFIG"`
	TrustedSubnet string `env:"TRUSTED_SUBNET"`
}

// FileConfig configuration file
type FileConfig struct {
	Addr            string `json:"server_address"`
	GRPCAddr        string `json:"grpc_address"`
	IsHTTPS         bool   `json:"enable_https"`
	FileStoragePath string `json:"file_storage_path"`
	DatabaseDSN     string `json:"database_dsn"`
	BaseURL         string `json:"base_url"`
	TrustedSubnet   string `json:"trusted_subnet"`
}

// New constructor a new instance of Configs
func New() (*Configs, error) {
	var cfg Configs

	fs := flag.NewFlagSet("shortener", flag.ContinueOnError)
	err := initFlag(&cfg, fs)
	if err != nil {
		return nil, err
	}

	err = initEnv(&cfg)
	if err != nil {
		return nil, err
	}

	err = initFile(&cfg, fs)
	if err != nil {
		return nil, err
	}

	initDefaultValue(&cfg)

	return &cfg, nil
}

func initFlag(cfg *Configs, fs *flag.FlagSet) error {
	fs.StringVar(&cfg.Server.Addr, "a", "", "Адрес запуска сервера localhost:8080")
	fs.StringVar(&cfg.GRPC.Addr, "g", "", "Адрес запуска grpc сервера localhost:3200")
	fs.StringVar(&cfg.BaseURL, "b", "", "Базовый URL localhost:8080")
	fs.StringVar(&cfg.Log.FlagLogLevel, "ll", "info", "log level")
	fs.StringVar(&cfg.File.FileStoragePath, "f", "", "Полное имя файла")
	fs.StringVar(&cfg.Cert.KeyFile, "fc", "key.pem", "Закрытый ключ")
	fs.StringVar(&cfg.Cert.CertFile, "fk", "cert.pem", "Подписанный центром сертификации, файл сертификата")
	fs.StringVar(&cfg.Database.DatabaseDSN, "d", "", "Строка с адресом подключения")
	fs.StringVar(&cfg.TrustedSubnet, "t", "", "Строковое представление бесклассовой адресации")
	fs.StringVar(&cfg.Config, "c", "", "Файл конфигурации")
	fs.BoolVar(&cfg.Server.IsHTTPS, "s", false, "Включения HTTPS в веб-сервере.")
	err := fs.Parse(os.Args[1:])
	if err != nil {
		return err
	}

	return nil
}

func initEnv(cfg *Configs) error {
	var err error
	if envServerAddress := os.Getenv("SERVER_ADDRESS"); len(envServerAddress) > 0 {
		cfg.Server.Addr = envServerAddress
	}

	if envGRPCServerAddress := os.Getenv("GRPC_ADDRESS"); len(envGRPCServerAddress) > 0 {
		cfg.GRPC.Addr = envGRPCServerAddress
	}

	if envBaseURL := os.Getenv("BASE_URL"); len(envBaseURL) > 0 {
		cfg.BaseURL = envBaseURL
	}

	if envIsHTTPS := os.Getenv("ENABLE_HTTPS"); len(envIsHTTPS) > 0 {
		cfg.Server.IsHTTPS, err = strconv.ParseBool(envIsHTTPS)
		if err != nil {
			return err
		}
	}

	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); len(envFileStoragePath) > 0 {
		cfg.File.FileStoragePath = envFileStoragePath
	}

	if envDatabaseDSN := os.Getenv("DATABASE_DSN"); len(envDatabaseDSN) > 0 {
		cfg.Database.DatabaseDSN = envDatabaseDSN
	}

	if envCertFile := os.Getenv("FILE_CERT"); len(envCertFile) > 0 {
		cfg.Cert.CertFile = envCertFile
	}

	if envKeyFile := os.Getenv("FILE_PRIVATE_KEY"); len(envKeyFile) > 0 {
		cfg.Cert.KeyFile = envKeyFile
	}

	if envFlagLogLevel := os.Getenv("FLAG_LOG_LEVEL"); len(envFlagLogLevel) > 0 {
		cfg.Log.FlagLogLevel = envFlagLogLevel
	}

	if envConfig := os.Getenv("CONFIG"); len(envConfig) > 0 {
		cfg.Config = envConfig
	}

	if envTrustedSubnet := os.Getenv("TRUSTED_SUBNET"); len(envTrustedSubnet) > 0 {
		cfg.TrustedSubnet = envTrustedSubnet
	}

	return nil
}

func initFile(cfg *Configs, fs *flag.FlagSet) error {

	if len(cfg.Config) > 0 {
		var fileConfig FileConfig
		file, err := os.Open(cfg.Config)
		if err != nil {
			return err
		}

		defer func() {
			_ = file.Close()
		}()

		decoder := json.NewDecoder(file)
		err = decoder.Decode(&fileConfig)
		if err != nil {
			return err
		}

		if len(cfg.Server.Addr) == 0 {
			cfg.Server.Addr = fileConfig.Addr
		}
		if len(cfg.GRPC.Addr) == 0 {
			cfg.GRPC.Addr = fileConfig.GRPCAddr
		}
		if len(cfg.BaseURL) == 0 {
			cfg.BaseURL = fileConfig.BaseURL
		}
		if len(cfg.Database.DatabaseDSN) == 0 {
			cfg.Database.DatabaseDSN = fileConfig.DatabaseDSN
		}

		if len(cfg.TrustedSubnet) == 0 {
			cfg.TrustedSubnet = fileConfig.TrustedSubnet
		}

		enabled := false
		fs.Visit(func(f *flag.Flag) {
			if f.Name == "s" {
				enabled = true
			}
		})
		envIsHTTPS := os.Getenv("ENABLE_HTTPS")

		if !cfg.Server.IsHTTPS && len(envIsHTTPS) == 0 && !enabled {
			cfg.Server.IsHTTPS = fileConfig.IsHTTPS
		}

		if len(cfg.File.FileStoragePath) == 0 {
			cfg.File.FileStoragePath = fileConfig.FileStoragePath
		}
	}

	return nil
}

func initDefaultValue(cfg *Configs) {
	cfg.Server.Addr = setValueString(cfg.Server.Addr, "localhost:8080")
	url := "localhost:8080"
	if cfg.Server.IsHTTPS {
		url = "https://" + url
	} else {
		url = "http://" + url
	}
	cfg.BaseURL = setValueString(cfg.BaseURL, url)
	cfg.File.FileStoragePath = setValueString(cfg.File.FileStoragePath, "/tmp/short-url-db.json")
}

func setValueString(value, defaultValue string) string {
	if len(value) > 0 {
		return value
	}
	return defaultValue
}
