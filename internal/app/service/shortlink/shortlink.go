package shortlink

import (
	"fmt"
	models "github.com/Orendev/shortener/internal/app/models/shortlink"
	repository "github.com/Orendev/shortener/internal/app/repository/shortlink"
	"github.com/Orendev/shortener/internal/app/service/filedb"
	"github.com/Orendev/shortener/internal/configs"
	"strings"
)

type Service struct {
	storage repository.ShortLinkRepository
	cfg     *configs.Configs
	fileDB  *filedb.FileDB
}

func NewService(storage repository.ShortLinkRepository, cfg *configs.Configs, fileDB *filedb.FileDB) *Service {

	return &Service{
		storage: storage,
		cfg:     cfg,
		fileDB:  fileDB,
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

	err = s.fileDB.Save(models.FileDB{
		OriginalURL: model.Link,
		ShortURL:    model.Code,
		UUID:        s.fileDB.ID(),
	})
	if err != nil {
		return model.Code, err
	}

	return code, nil
}
