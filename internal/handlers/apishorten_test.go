package handlers_test

import (
	"bytes"
	"github.com/Orendev/shortener/internal/handlers"
	"github.com/Orendev/shortener/internal/middlewares"
	mockStore "github.com/Orendev/shortener/internal/repository/mock"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_DeleteAPIUserUrls(t *testing.T) {

	// создадим конроллер моков и экземпляр мок-хранилища
	ctrl := gomock.NewController(t)
	s := mockStore.NewMockStorage(ctrl)

	// определим, какой результат будем получать от «хранилища»
	// установим условие: при любом вызове метода Save возвращать uuid без ошибки
	s.EXPECT().
		DeleteFlagBatch(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)

	// создадим экземпляр приложения и передадим ему «хранилище»
	h := handlers.NewHandler(s, "http://localhost")

	r := chi.NewRouter()
	r.Use(middlewares.Auth)
	r.Delete("/api/user/urls", h.DeleteAPIUserUrls)

	srv := httptest.NewServer(r)
	defer srv.Close()

	type want struct {
		expectedCode int
		contentType  string
	}
	tests := []struct {
		name   string
		method string
		body   string // добавляем тело запроса в табличные тесты
		want   want
	}{
		{
			name:   "method delete success DeleteAPIUserUrls",
			method: http.MethodDelete,
			body:   `["zIlkHekB"]`,
			want: want{
				expectedCode: http.StatusAccepted,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bodyReader io.Reader

			if len(tt.body) > 0 {
				jsonBody := []byte(tt.body)
				bodyReader = bytes.NewReader(jsonBody)
			}

			req, err := http.NewRequest(tt.method, srv.URL+"/api/user/urls", bodyReader)
			require.NoError(t, err)

			if len(tt.body) > 0 {
				req.Header.Set("Content-Type", "application/json")
			}

			resp, err := srv.Client().Do(req)
			require.NoError(t, err)

			defer func() {
				err = resp.Body.Close()
				if err != nil {
					require.NoError(t, err)
				}
			}()

			assert.NoError(t, err, "error making HTTP request")
			assert.Equal(t, tt.want.expectedCode, resp.StatusCode, "code didn't match expected")

		})
	}
}
