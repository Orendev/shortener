package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	models "github.com/Orendev/shortener/internal/app/models/shortlink"
	"github.com/Orendev/shortener/internal/configs"
	"github.com/google/uuid"
	"log"
	"os"
	"strings"
)

type FileDB struct {
	data map[string]models.ShortLink
	cfg  *configs.Configs
}

func NewFileDB(cfg *configs.Configs) (*FileDB, error) {
	return &FileDB{
		cfg:  cfg,
		data: cfg.Memory,
	}, nil
}

func (f *FileDB) ID() string {
	id := uuid.New()
	return id.String()
}

// Save сохраняет данные в файле FileStoragePath.
func (f *FileDB) Save(fileDB models.FileDB) error {

	file, err := os.OpenFile(f.cfg.FileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatalf("error while closing file: %s", err)
		}
	}()

	// сериализуем структуру в JSON формат
	writeData, err := json.Marshal(fileDB)
	if err != nil {
		return err
	}
	_, err = file.Write(append(writeData, '\n'))
	if err != nil {
		return err
	}
	return nil
}

// Load Прочитаем данные из файла FileStoragePath
func (f *FileDB) Load() error {

	file, err := os.OpenFile(f.cfg.FileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	scan := bufio.NewScanner(file)

	for {
		if !scan.Scan() {
			break
		}
		fileDB := models.FileDB{}
		data := scan.Bytes()
		err = json.Unmarshal(data, &fileDB)

		if err != nil {
			log.Fatal(err)
		}

		f.data[fileDB.ShortURL] = models.ShortLink{
			Link:   fileDB.OriginalURL,
			Code:   fileDB.ShortURL,
			Result: fmt.Sprintf("%s/%s", strings.TrimPrefix(f.cfg.BaseURL, "/"), fileDB.ShortURL),
		}

	}
	return nil
}

// Remove Удалим файл FileStoragePath
func (f *FileDB) Remove() error {
	return os.Remove(f.cfg.FileStoragePath)
}
