package repository

import (
	"testing"

	"github.com/Orendev/shortener/internal/config"
	"github.com/Orendev/shortener/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestFileDB_Remove(t *testing.T) {
	type args struct {
		models map[string]models.ShortLink
	}
	id := uuid.New().String()

	model := models.ShortLink{
		UUID:        id,
		OriginalURL: "http://yandex.ru",
		ShortURL:    "http://localhost:8080/4rSPg8ap",
	}

	tests := []struct {
		name    string
		cfg     *config.Configs
		args    args
		wantErr bool
	}{
		{
			name: "test File Remove",
			args: args{
				models: map[string]models.ShortLink{
					"4rSPg8ap": model,
				},
			},
			cfg: &config.Configs{
				Server: config.Server{
					Host: "",
					Port: "8080",
				},
				BaseURL: "http://localhost:8080",
				Memory:  map[string]models.ShortLink{},
				File: config.File{
					FileStoragePath: "/tmp/test-short-url-file.json",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := File{
				data: tt.cfg.Memory,
				cfg:  tt.cfg,
			}

			if err := f.Save(tt.args.models); (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := f.Remove(); (err != nil) != tt.wantErr {
				t.Errorf("Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFileDB_Save(t *testing.T) {

	type args struct {
		models map[string]models.ShortLink
	}
	id := uuid.New().String()
	model := models.ShortLink{
		UUID:        id,
		OriginalURL: "http://yandex.ru",
		ShortURL:    "http://localhost:8080/4rSPg8ap",
	}
	tests := []struct {
		name    string
		cfg     *config.Configs
		args    args
		wantErr bool
	}{
		{
			name: "test DB Save",
			args: args{
				models: map[string]models.ShortLink{
					"4rSPg8ap": model,
				},
			},
			cfg: &config.Configs{
				Server: config.Server{
					Host: "",
					Port: "8080",
				},
				BaseURL: "http://localhost:8080",
				Memory:  map[string]models.ShortLink{},
				File: config.File{
					FileStoragePath: "/tmp/test-short-url-file.json",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := File{
				data: tt.cfg.Memory,
				cfg:  tt.cfg,
			}
			if err := f.Save(tt.args.models); (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}

			err := f.Remove()
			if err != nil {
				require.NoError(t, err)
			}
		})
	}
}
