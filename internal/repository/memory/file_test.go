package memory

import (
	"reflect"
	"testing"

	"github.com/Orendev/shortener/internal/config"
	"github.com/Orendev/shortener/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestFile_Remove(t *testing.T) {
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
				Server:  config.Server{},
				BaseURL: "http://localhost:8080",
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
				filePath: tt.cfg.File.FileStoragePath,
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

func TestFile_Save(t *testing.T) {

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
			name: "test Save",
			args: args{
				models: map[string]models.ShortLink{
					"4rSPg8ap": model,
				},
			},
			cfg: &config.Configs{
				Server:  config.Server{},
				BaseURL: "http://localhost:8080",
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
				filePath: tt.cfg.File.FileStoragePath,
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

func TestFile_Data(t *testing.T) {
	id := uuid.New().String()

	model := models.ShortLink{
		UUID:        id,
		Code:        "4rSPg8ap",
		OriginalURL: "http://yandex.ru",
		ShortURL:    "http://localhost:8080/4rSPg8ap",
	}

	type fields struct {
		filePath string
	}
	tests := []struct {
		name    string
		fields  fields
		want    map[string]models.ShortLink
		wantErr bool
	}{
		{
			name: "test file Data",
			fields: fields{
				filePath: "/tmp/test-short-url-file.json",
			},
			want: map[string]models.ShortLink{
				"4rSPg8ap": model,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &File{
				filePath: tt.fields.filePath,
			}
			if err := f.Save(tt.want); (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}

			got, err := f.Data()
			if (err != nil) != tt.wantErr {
				t.Errorf("Data() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Data() got = %v, want %v", got, tt.want)
			}

			err = f.Remove()
			if err != nil {
				require.NoError(t, err)
			}
		})
	}
}
