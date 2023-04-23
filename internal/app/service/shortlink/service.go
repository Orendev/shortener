package shortlink

import (
	"fmt"
	models "github.com/Orendev/shortener/internal/app/models/shortlink"
	repository "github.com/Orendev/shortener/internal/app/repository/shortlink"
	"github.com/Orendev/shortener/internal/configs"
	"strings"
)

type Service struct {
	storage repository.ShortLinkRepository
	cfg     *configs.Configs
}

func NewService(storage repository.ShortLinkRepository, cfg *configs.Configs) *Service {
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
	return s.storage.Add(model)
}
