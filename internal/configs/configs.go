package configs

import "github.com/Orendev/shortener/internal/app/server"

type Configs struct {
	Host string
	Port string
}

func (cfg Configs) Server() (*server.Config, error) {
	return &server.Config{
		Port: "8080",
	}, nil
}

func New() (*Configs, error) {
	return &Configs{}, nil
}
