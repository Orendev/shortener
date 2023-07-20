package memory

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/Orendev/shortener/internal/logger"
	"github.com/Orendev/shortener/internal/models"
)

type File struct {
	filePath string
}

func NewFile(filePath string) *File {
	return &File{
		filePath: filePath,
	}
}

// Save сохраняет данные в файле FileStoragePath.
func (f *File) Save(models map[string]models.ShortLink) error {

	file, err := os.OpenFile(f.filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			logger.Log.Sugar().Errorf("error while closing file: %s", err)
		}
	}()

	for _, model := range models {
		// сериализуем структуру в JSON формат
		writeData, err := json.Marshal(model)
		if err != nil {
			return err
		}
		_, err = file.Write(append(writeData, '\n'))
		if err != nil {
			return err
		}
	}

	return nil
}

// Remove Удалим файл FileStoragePath
func (f *File) Remove() error {
	return os.Remove(f.filePath)
}

// Data Прочитаем данные из файла FileStoragePath
func (f *File) Data() (map[string]models.ShortLink, error) {
	file, err := os.OpenFile(f.filePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := file.Close(); err != nil {
			logger.Log.Sugar().Errorf("error when closing a file while reading: %s", err)
		}
	}()

	scan := bufio.NewScanner(file)
	data := make(map[string]models.ShortLink)

	for {
		if !scan.Scan() {
			break
		}
		model := models.ShortLink{}

		err = json.Unmarshal(scan.Bytes(), &model)

		if err != nil {
			return nil, err
		}

		data[model.Code] = model

	}

	return data, nil
}
