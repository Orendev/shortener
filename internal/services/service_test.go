package services

import (
	"context"
	"github.com/Orendev/shortener/internal/models"
	"github.com/Orendev/shortener/internal/random"
	"github.com/Orendev/shortener/internal/storage"
	"github.com/Orendev/shortener/internal/storage/mock"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestService_Add(t *testing.T) {
	// создадим конроллер моков и экземпляр мок-хранилища
	ctrl := gomock.NewController(t)
	s := mock.NewMockShortLinkStorage(ctrl)

	code := random.Strn(8)

	type args struct {
		model models.ShortLink
	}

	tests := []struct {
		name string // добавим название тестов
		args args
		want string
	}{
		{
			name: "positive test #1 method Add storage",
			args: args{
				model: models.ShortLink{
					UUID:        uuid.New().String(),
					Code:        code,
					ShortURL:    "http://localhost/" + code,
					OriginalURL: "https://practicum.yandex.ru/",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// установим условие: при любом вызове метода Add возвращать uuid без ошибки
			s.EXPECT().
				Add(gomock.Any(), gomock.Any()).
				Return(tt.args.model.UUID, nil)

			service := NewService(s)

			val, err := service.Add(context.Background(), &tt.args.model)
			// и проверяем возвращаемые значения
			require.NoError(t, err)
			require.Equal(t, val, tt.args.model.UUID)
		})
	}
}

func TestService_Close(t *testing.T) {

	// создадим конроллер моков и экземпляр мок-хранилища
	ctrl := gomock.NewController(t)
	s := mock.NewMockShortLinkStorage(ctrl)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "positive test #1 method Close storage",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// установим условие: при любом вызове метода Add возвращать uuid без ошибки
			s.EXPECT().
				Close().
				Return(nil)

			service := NewService(s)

			if err := service.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_GetByCode(t *testing.T) {
	// создадим конроллер моков и экземпляр мок-хранилища
	ctrl := gomock.NewController(t)
	s := mock.NewMockShortLinkStorage(ctrl)

	code := random.Strn(8)
	model := models.ShortLink{
		UUID:        uuid.New().String(),
		Code:        code,
		ShortURL:    "http://localhost/" + code,
		OriginalURL: "https://practicum.yandex.ru/",
	}

	type fields struct {
		storage storage.ShortLinkStorage
	}
	type args struct {
		ctx  context.Context
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
			name: "positive test #1 method GetByCode storage",
			want: &model,
			args: args{
				ctx:  context.Background(),
				code: code,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// установим условие: при любом вызове метода Add возвращать uuid без ошибки

			s.EXPECT().
				GetByCode(gomock.Any(), gomock.Any()).
				Return(&model, nil)

			service := NewService(s)

			val, err := service.GetByCode(tt.args.ctx, tt.args.code)

			// и проверяем возвращаемые значения
			require.NoError(t, err)
			require.Equal(t, val.UUID, model.UUID)

			if !reflect.DeepEqual(val, &model) {
				t.Errorf("GetByCode() got = %v, want %v", val, model)
			}

			//s := &Service{
			//	storage: tt.fields.storage,
			//}
			//got, err := s.GetByCode(tt.args.ctx, tt.args.code)
			//if (err != nil) != tt.wantErr {
			//	t.Errorf("GetByCode() error = %v, wantErr %v", err, tt.wantErr)
			//	return
			//}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("GetByCode() got = %v, want %v", got, tt.want)
			//}
		})
	}
}

func TestService_Ping(t *testing.T) {
	// создадим конроллер моков и экземпляр мок-хранилища
	ctrl := gomock.NewController(t)
	s := mock.NewMockShortLinkStorage(ctrl)

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "positive test #1 method Ping storage",
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			s.EXPECT().
				Ping(gomock.Any()).
				Return(nil)

			service := NewService(s)

			// и проверяем возвращаемые значения
			if err := service.Ping(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Ping() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
