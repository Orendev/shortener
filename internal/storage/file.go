package storage

import (
	"bufio"
	"encoding/json"
	"github.com/Orendev/shortener/internal/configs"
	"github.com/Orendev/shortener/internal/models"
	"log"
	"os"
)

type File struct {
	data map[string]models.ShortLink
	cfg  *configs.Configs
}

func NewFile(cfg *configs.Configs) (*File, error) {
	//Прочитаем данные из файла FileStoragePath
	file, err := os.OpenFile(cfg.FileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)
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
		model := models.ShortLink{}
		data := scan.Bytes()

		err = json.Unmarshal(data, &model)

		if err != nil {
			log.Fatal(err)
		}

		cfg.Memory[model.Code] = model

	}

	return &File{
		cfg:  cfg,
		data: cfg.Memory,
	}, nil
}

// Save сохраняет данные в файле FileStoragePath.
func (f *File) Save(model models.ShortLink) error {

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
	writeData, err := json.Marshal(model)
	if err != nil {
		return err
	}
	_, err = file.Write(append(writeData, '\n'))
	if err != nil {
		return err
	}
	return nil
}

// Remove Удалим файл FileStoragePath
func (f *File) Remove() error {
	return os.Remove(f.cfg.FileStoragePath)
}
