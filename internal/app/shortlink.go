package app

import (
	"fmt"
	"github.com/Orendev/shortener/internal/configs"
	models "github.com/Orendev/shortener/internal/models/shortlink"
	"github.com/Orendev/shortener/internal/storage"
	"strings"
)

type Service struct {
	storage storage.ShortLinkStorage
	cfg     *configs.Configs
}

func NewService(storage storage.ShortLinkStorage, cfg *configs.Configs) *Service {

	return &Service{
		storage: storage,
		cfg:     cfg,
	}
}

func (s *Service) Get(code string) (*models.ShortLink, error) {
	return s.storage.Get(code)
}

func (s *Service) Add(model *models.ShortLink) (string, error) {
	model.Result = fmt.Sprintf("%s/%s", strings.TrimPrefix(s.cfg.BaseURL, "/"), model.Code)

	code, err := s.storage.Add(model)
	if err != nil {
		return model.Code, err
	}

	return code, nil
}
