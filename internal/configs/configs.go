package configs

import (
	"flag"
)

var options Configs = Configs{}

type Configs struct {
	Addr    string
	Host    string
	Port    string
	BaseURL string
}

func init() {
	flag.StringVar(&options.Addr, "a", "", "Адрес запуска сервера localhost:8080")
	flag.StringVar(&options.BaseURL, "b", "http://localhost:8080", "Базовый URL http://localhost:8080")
	flag.StringVar(&options.Port, "p", "8080", "Порт 8080")
}

func New() (*Configs, error) {
	flag.Parse()
	return &options, nil
}
