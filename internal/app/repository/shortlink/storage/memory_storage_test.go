package storage

import (
	"github.com/Orendev/shortener/internal/app/repository/shortlink/model"
	"github.com/Orendev/shortener/internal/configs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestMemoryStorage_Get(t *testing.T) {
	type fields struct {
		data map[string]model.ShortLink
		cfg  *configs.Configs
	}
	type args struct {
		code string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.ShortLink
		wantErr bool
	}{
		{
			name: "positive test #1 memory storage",
			args: args{
				code: "test",
			},
			fields: fields{
				data: map[string]model.ShortLink{
					"test": {
						Code: "test",
						Link: "localhost",
					},
				},
			},
			want: &model.ShortLink{
				Code: "test",
				Link: "localhost",
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
			got, err := s.Get(tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}

			assert.Equal(t, got.Code, tt.want.Code)
			assert.Equal(t, got.Link, tt.want.Link)
		})
	}
}

func TestMemoryStorage_Add(t *testing.T) {
	type fields struct {
		data map[string]model.ShortLink
		cfg  *configs.Configs
	}
	type args struct {
		shortLink model.ShortLink
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
				shortLink: model.ShortLink{
					Code: "test",
					Link: "localhost",
				},
			},
			fields: fields{
				data: map[string]model.ShortLink{},
				cfg: &configs.Configs{
					Host:    "",
					Port:    "8080",
					BaseURL: "http://localhost:8080",
					Memory:  map[string]model.ShortLink{},
				},
			},
			want: "http://localhost:8080/test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemoryStorage{
				data: tt.fields.data,
				cfg:  tt.fields.cfg,
			}
			got, err := s.Add(tt.args.shortLink)
			require.NoError(t, err)

			assert.Equalf(t, tt.want, got, "Add(%v)", tt.args.shortLink)
		})
	}
}