package shortlink

import (
	"fmt"
	model "github.com/Orendev/shortener/internal/app/models/shortlink"
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

func (s *Service) Get(code string) (*model.ShortLink, error) {
	return s.storage.Get(code)
}

func (s *Service) Add(model model.ShortLink) (string, error) {
	code, err := s.storage.Add(model)
	url := fmt.Sprintf("%s/%s", strings.TrimPrefix(s.cfg.BaseURL, "/"), code)
	fmt.Println(url)
	return url, err
}
