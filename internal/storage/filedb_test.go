package storage

import (
	"github.com/Orendev/shortener/internal/configs"
	models "github.com/Orendev/shortener/internal/models/shortlink"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFileDB_ID(t *testing.T) {
	type args struct {
		fileDB models.FileDB
	}
	id := uuid.New().String()
	tests := []struct {
		name    string
		cfg     *configs.Configs
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test fileDB ID",
			args: args{
				fileDB: models.FileDB{
					UUID:        id,
					OriginalURL: "http://yandex.ru",
					ShortURL:    "4rSPg8ap",
				},
			},
			cfg: &configs.Configs{
				Host:            "",
				Port:            "8080",
				BaseURL:         "http://localhost:8080",
				Memory:          map[string]models.ShortLink{},
				FileStoragePath: "/tmp/test-short-url-db.json",
			},
			want:    id,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := NewFileDB(tt.cfg)
			require.NoError(t, err)

			_, err = uuid.Parse(f.ID())
			require.NoError(t, err)

		})
	}
}

func TestFileDB_Load(t *testing.T) {
	type args struct {
		fileDB models.FileDB
	}
	id := uuid.New().String()
	tests := []struct {
		name    string
		cfg     *configs.Configs
		args    args
		wantErr bool
	}{
		{
			name: "test FileDB Load",
			args: args{
				fileDB: models.FileDB{
					UUID:        id,
					OriginalURL: "http://yandex.ru",
					ShortURL:    "4rSPg8ap",
				},
			},
			cfg: &configs.Configs{
				Host:            "",
				Port:            "8080",
				BaseURL:         "http://localhost:8080",
				Memory:          map[string]models.ShortLink{},
				FileStoragePath: "/tmp/test-short-url-db.json",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FileDB{
				data: tt.cfg.Memory,
				cfg:  tt.cfg,
			}

			if err := f.Save(tt.args.fileDB); (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := f.Load(); (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, tt.cfg.Memory[tt.args.fileDB.ShortURL].Code, tt.args.fileDB.ShortURL)

			err := f.Remove()
			require.NoError(t, err)
		})
	}
}

func TestFileDB_Remove(t *testing.T) {
	type args struct {
		fileDB models.FileDB
	}
	id := uuid.New().String()
	tests := []struct {
		name    string
		cfg     *configs.Configs
		args    args
		wantErr bool
	}{
		{
			name: "test FileDB Remove",
			args: args{
				fileDB: models.FileDB{
					UUID:        id,
					OriginalURL: "http://yandex.ru",
					ShortURL:    "4rSPg8ap",
				},
			},
			cfg: &configs.Configs{
				Host:            "",
				Port:            "8080",
				BaseURL:         "http://localhost:8080",
				Memory:          map[string]models.ShortLink{},
				FileStoragePath: "/tmp/test-short-url-db.json",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FileDB{
				data: tt.cfg.Memory,
				cfg:  tt.cfg,
			}

			if err := f.Save(tt.args.fileDB); (err != nil) != tt.wantErr {
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
		fileDB models.FileDB
	}
	id := uuid.New().String()
	tests := []struct {
		name    string
		cfg     *configs.Configs
		args    args
		wantErr bool
	}{
		{
			name: "test DB Save",
			args: args{
				fileDB: models.FileDB{
					UUID:        id,
					OriginalURL: "http://yandex.ru",
					ShortURL:    "4rSPg8ap",
				},
			},
			cfg: &configs.Configs{
				Host:            "",
				Port:            "8080",
				BaseURL:         "http://localhost:8080",
				Memory:          map[string]models.ShortLink{},
				FileStoragePath: "/tmp/test-short-url-db.json",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FileDB{
				data: tt.cfg.Memory,
				cfg:  tt.cfg,
			}
			if err := f.Save(tt.args.fileDB); (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}

			err := f.Remove()
			if err != nil {
				require.NoError(t, err)
			}
		})
	}
}
