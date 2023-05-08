package storage

import (
	"github.com/Orendev/shortener/internal/configs"
	"github.com/Orendev/shortener/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestMemoryStorage_Get(t *testing.T) {
	type fields struct {
		data map[string]models.ShortLink
		cfg  *configs.Configs
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
						OriginalUrl: "localhost",
					},
				},
			},
			want: &models.ShortLink{
				Code:        "test",
				OriginalUrl: "localhost",
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
			got, err := s.GetByCode(tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}

			assert.Equal(t, got.Code, tt.want.Code)
			assert.Equal(t, got.OriginalUrl, tt.want.OriginalUrl)
		})
	}
}

func TestMemoryStorage_Add(t *testing.T) {
	id := uuid.New().String()
	type fields struct {
		data map[string]models.ShortLink
		cfg  *configs.Configs
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
					OriginalUrl: "localhost",
				},
			},
			fields: fields{
				data: map[string]models.ShortLink{},
				cfg: &configs.Configs{
					Host:            "",
					Port:            "8080",
					BaseURL:         "http://localhost:8080",
					Memory:          map[string]models.ShortLink{},
					FileStoragePath: "/tmp/test-short-url-db.json",
				},
			},
			want: id,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileDB, err := NewFileDB(tt.fields.cfg)
			if err != nil {
				require.NoError(t, err)
			}

			s := &MemoryStorage{
				data: tt.fields.data,
				cfg:  tt.fields.cfg,
				db:   fileDB,
			}
			got, err := s.Add(&tt.args.shortLink)
			require.NoError(t, err)

			assert.Equalf(t, tt.want, got, "Add(%v)", tt.args.shortLink)

			err = s.db.Remove()
			require.NoError(t, err)
		})
	}
}
