package storage

import (
	"context"
	"github.com/Orendev/shortener/internal/config"
	"github.com/Orendev/shortener/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestMemoryStorage_GetByCode(t *testing.T) {
	type fields struct {
		data map[string]models.ShortLink
		cfg  *config.Configs
	}
	type args struct {
		code string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.ShortLink
		wantErr bool
	}{
		{
			name: "positive test #1 memory storage",
			args: args{
				code: "test",
			},
			fields: fields{
				data: map[string]models.ShortLink{
					"test": {
						Code:        "test",
						OriginalURL: "localhost",
					},
				},
			},
			want: &models.ShortLink{
				Code:        "test",
				OriginalURL: "localhost",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemoryStorage{
				data: tt.fields.data,
				cfg:  tt.fields.cfg,
			}
			got, err := s.GetByCode(context.Background(), tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}

			assert.Equal(t, got.Code, tt.want.Code)
			assert.Equal(t, got.OriginalURL, tt.want.OriginalURL)
		})
	}
}

func TestMemoryStorage_Add(t *testing.T) {
	id := uuid.New().String()
	type fields struct {
		data map[string]models.ShortLink
		cfg  *config.Configs
	}
	type args struct {
		shortLink models.ShortLink
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "positive test #2 memory storage",
			args: args{
				shortLink: models.ShortLink{
					UUID:        id,
					Code:        "test",
					OriginalURL: "localhost",
				},
			},
			fields: fields{
				data: map[string]models.ShortLink{},
				cfg: &config.Configs{
					Host:            "",
					Port:            "8080",
					BaseURL:         "http://localhost:8080",
					Memory:          map[string]models.ShortLink{},
					FileStoragePath: "/tmp/test-short-url-file.json",
					DatabaseDSN:     "host=localhost user=shortener password=secret dbname=shortener sslmode=disable",
				},
			},
			want: id,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := NewFile(tt.fields.cfg)
			require.NoError(t, err)

			s := &MemoryStorage{
				data: tt.fields.data,
				cfg:  tt.fields.cfg,
				file: file,
			}
			err = s.Save(context.Background(), tt.args.shortLink)
			require.NoError(t, err)

			defer func() {
				err = file.Remove()
				if err != nil {
					require.NoError(t, err)
				}
			}()

		})
	}
}
