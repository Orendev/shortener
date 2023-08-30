package handlers_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Orendev/shortener/internal/handlers"
	"github.com/Orendev/shortener/internal/middlewares"
	"github.com/Orendev/shortener/internal/models"
	"github.com/Orendev/shortener/internal/random"
	"github.com/Orendev/shortener/internal/repository/mock"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_GetShorten(t *testing.T) {

	// создадим конроллер моков и экземпляр мок-хранилища
	ctrl := gomock.NewController(t)
	s := mockStore.NewMockStorage(ctrl)

	// определим, какой результат будем получать от «хранилища»
	code := random.Strn(8)
	// определим, какой результат будем получать от «хранилища»
	model := models.ShortLink{
		UUID:        uuid.New().String(),
		Code:        code,
		ShortURL:    "http://localhost/" + code,
		OriginalURL: "https://practicum.yandex.ru/",
	}

	// установим условие: при любом вызове метода Save возвращать uuid без ошибки
	s.EXPECT().
		GetByCode(gomock.Any(), gomock.Any()).
		Return(&model, nil)

	// создадим экземпляр приложения и передадим ему «хранилище»
	h := handlers.NewHandler(s, "http://localhost", "192.168.1.0/24")

	srv := httptest.NewServer(http.HandlerFunc(h.GetShorten))
	defer srv.Close()

	type want struct {
		expectedCode int
		contentType  string
	}

	tests := []struct {
		name   string
		method string
		want   want
	}{
		{
			name:   "positive test #1 ShortLink",
			method: http.MethodGet,
			want: want{
				expectedCode: http.StatusTemporaryRedirect,
				contentType:  "text/plain",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, srv.URL+"/"+code, nil)
			require.NoError(t, err)

			//Similarly, RoundTrip should not attempt to
			//handle higher-level protocol details such as redirects,
			resp, err := srv.Client().Transport.RoundTrip(req)
			require.NoError(t, err)

			defer func() {
				err := resp.Body.Close()
				if err != nil {
					require.NoError(t, err)
				}
			}()

			assert.Equal(t, model.OriginalURL, resp.Header.Get("Location"))
			assert.Equal(t, tt.want.expectedCode, resp.StatusCode, "code didn't match expected")

		})
	}
}

func TestHandler_PostShorten(t *testing.T) {

	// создадим конроллер моков и экземпляр мок-хранилища
	ctrl := gomock.NewController(t)
	s := mockStore.NewMockStorage(ctrl)

	// определим, какой результат будем получать от «хранилища»
	// установим условие: при любом вызове метода Save возвращать uuid без ошибки
	s.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(nil)

	// создадим экземпляр приложения и передадим ему «хранилище»
	h := handlers.NewHandler(s, "http://localhost", "192.168.1.0/24")

	r := chi.NewRouter()
	r.Use(middlewares.Auth)
	r.Post("/", h.PostShorten)

	srv := httptest.NewServer(r)
	defer srv.Close()

	type want struct {
		expectedCode int
		expectedBody string
		contentType  string
	}

	tests := []struct {
		name   string
		method string
		body   string // добавляем тело запроса в табличные тесты
		want   want
	}{
		{
			name:   "positive test #1 ShortLinkAdd",
			method: http.MethodPost,
			body:   `https://google.com`,
			want: want{
				expectedCode: http.StatusCreated,
				contentType:  "text/plain",
				expectedBody: "http://localhost/.*",
			},
		},
		{
			name:   "negative test #2 ShortLinkAdd",
			method: http.MethodPost,
			want: want{
				expectedCode: http.StatusBadRequest,
				contentType:  "text/plain; charset=utf-8",
				expectedBody: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var bodyReader io.Reader
			if len(tt.body) > 0 {
				bodyReader = bytes.NewReader([]byte(tt.body))
			}

			req, err := http.NewRequest(tt.method, srv.URL+"/", bodyReader)
			require.NoError(t, err)

			resp, err := srv.Client().Do(req)
			require.NoError(t, err)

			defer func() {
				err := resp.Body.Close()
				if err != nil {
					require.NoError(t, err)
				}
			}()

			//assert.NoError(t, err, "error making HTTP request")
			assert.Equal(t, tt.want.expectedCode, resp.StatusCode, "code didn't match expected")

			// проверяем корректность полученного тела ответа, если мы его ожидаем
			if tt.want.expectedBody != "" {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					require.NoError(t, err)
				}

				assert.Regexp(t, tt.want.expectedBody, string(body))
			}

		})
	}
}

func TestHandler_PostAPIShorten(t *testing.T) {

	// создадим конроллер моков и экземпляр мок-хранилища
	ctrl := gomock.NewController(t)
	s := mockStore.NewMockStorage(ctrl)

	// определим, какой результат будем получать от «хранилища»
	// установим условие: при любом вызове метода Save возвращать uuid без ошибки
	s.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(nil)

	// создадим экземпляр приложения и передадим ему «хранилище»
	h := handlers.NewHandler(s, "http://localhost", "192.168.1.0/24")
	r := chi.NewRouter()
	r.Use(middlewares.Auth)
	r.Post("/api/shorten", h.PostAPIShorten)
	srv := httptest.NewServer(r)
	defer srv.Close()

	type want struct {
		contentType  string
		expectedCode int
		expectedBody string
	}

	tests := []struct {
		name   string // добавляем название тестов
		method string
		body   string // добавляем тело запроса в табличные тесты
		want   want
	}{
		{
			name:   "method_get",
			method: http.MethodGet,
			want: want{
				expectedCode: http.StatusMethodNotAllowed,
				expectedBody: "",
			},
		},
		{
			name:   "method_put",
			method: http.MethodPut,
			want: want{
				expectedCode: http.StatusMethodNotAllowed,
				expectedBody: "",
			},
		},
		{
			name:   "method_delete",
			method: http.MethodDelete,
			want: want{
				expectedCode: http.StatusMethodNotAllowed,
				expectedBody: "",
			},
		},
		{
			name:   "method_post_without_body",
			method: http.MethodPost,
			want: want{
				expectedCode: http.StatusInternalServerError,
				expectedBody: "",
			},
		},
		{
			name:   "method_post_success",
			method: http.MethodPost,
			body:   `{"url":"https://practicum.yandex.ru/"}`,
			want: want{
				expectedCode: http.StatusCreated,
				expectedBody: `{"result":"http://localhost/.*"}`,
			},
		},
		{
			name:   "method_post_bad_request",
			method: http.MethodPost,
			body:   `{"url":""}`,
			want: want{
				expectedCode: http.StatusBadRequest,
				expectedBody: "",
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

			req, err := http.NewRequest(tt.method, srv.URL+"/api/shorten", bodyReader)
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

			// проверяем корректность полученного тела ответа, если мы его ожидаем
			if tt.want.expectedBody != "" {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					require.NoError(t, err)
				}
				assert.Regexp(t, tt.want.expectedBody, string(body))
			}
		})
	}
}

func TestHandler_PostAPIShortenBatch(t *testing.T) {

	// создадим конроллер моков и экземпляр мок-хранилища
	ctrl := gomock.NewController(t)
	s := mockStore.NewMockStorage(ctrl)

	// определим, какой результат будем получать от «хранилища»
	code := random.Strn(8)
	id := uuid.New().String()
	// определим, какой результат будем получать от «хранилища»
	model := models.ShortLink{
		UUID:        id,
		Code:        code,
		ShortURL:    "http://localhost/" + code,
		OriginalURL: "https://practicum.yandex.ru/",
	}

	// определим, какой результат будем получать от «хранилища»
	// установим условие: при любом вызове метода Save возвращать uuid без ошибки
	s.EXPECT().
		InsertBatch(gomock.Any(), gomock.Any()).
		Return(nil)

	s.EXPECT().
		UpdateBatch(gomock.Any(), gomock.Any()).
		Return(nil)

	s.EXPECT().
		GetByID(gomock.Any(), gomock.Any()).
		Return(&model, nil)

	// создадим экземпляр приложения и передадим ему «хранилище»
	h := handlers.NewHandler(s, "http://localhost", "192.168.1.0/24")

	r := chi.NewRouter()
	r.Use(middlewares.Auth)
	r.Post("/api/shorten/batch", h.PostAPIShortenBatch)

	srv := httptest.NewServer(r)
	defer srv.Close()

	type want struct {
		expectedCode int
		expectedBody string
		contentType  string
	}

	tests := []struct {
		name   string
		method string
		body   string // добавляем тело запроса в табличные тесты
		want   want
	}{
		{
			name:   "method post success ShortenBatchInsert",
			method: http.MethodPost,
			body:   `[{"correlation_id": "e8cd3fd9-d161-4d47-9337-e09eb6ec0124", "original_url":"https://practicum.yandex.ru/"}]`,
			want: want{
				expectedCode: http.StatusCreated,
				expectedBody: `[{"correlation_id": "e8cd3fd9-d161-4d47-9337-e09eb6ec0124", "short_url":"http://localhost/.*"}]`,
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

			req, err := http.NewRequest(tt.method, srv.URL+"/api/shorten/batch", bodyReader)
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

			// проверяем корректность полученного тела ответа, если мы его ожидаем
			if tt.want.expectedBody != "" {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					require.NoError(t, err)
				}
				assert.Regexp(t, tt.want.expectedBody, string(body))
			}
		})
	}
}

func TestHandler_GetAPIUserUrls(t *testing.T) {
	// создадим конроллер моков и экземпляр мок-хранилища
	ctrl := gomock.NewController(t)
	s := mockStore.NewMockStorage(ctrl)

	// определим, какой результат будем получать от «хранилища»
	code := random.Strn(8)
	id := uuid.New().String()
	userID := uuid.New().String()
	// определим, какой результат будем получать от «хранилища»
	model := models.ShortLink{
		UUID:        id,
		UserID:      userID,
		Code:        code,
		ShortURL:    "http://localhost/" + code,
		OriginalURL: "https://practicum.yandex.ru/",
	}
	shortLinks := make([]models.ShortLink, 0)

	shortLinks = append(shortLinks, model)
	// определим, какой результат будем получать от «хранилища»
	s.EXPECT().
		ShortLinksByUserID(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(shortLinks, nil)

	// создадим экземпляр приложения и передадим ему «хранилище»
	h := handlers.NewHandler(s, "http://localhost", "192.168.1.0/24")

	r := chi.NewRouter()
	r.Use(middlewares.Auth)
	r.Get("/api/user/urls", h.GetAPIUserUrls)

	srv := httptest.NewServer(r)
	defer srv.Close()

	type want struct {
		expectedCode int
		expectedBody string
		contentType  string
	}

	type args struct {
		userID string
	}

	tests := []struct {
		name   string
		method string
		args   args
		want   want
	}{
		{
			name:   "method get success",
			method: http.MethodGet,
			args: args{
				userID: userID,
			},
			want: want{
				contentType:  "application/json",
				expectedCode: http.StatusOK,
				expectedBody: `[{"original_url": "https://practicum.yandex.ru/", "short_url":"http://localhost/.*"}]`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bodyReader io.Reader
			req, err := http.NewRequest(tt.method, srv.URL+"/api/user/urls", bodyReader)
			require.NoError(t, err)

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

			// проверяем корректность полученного тела ответа, если мы его ожидаем
			if tt.want.expectedBody != "" {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					require.NoError(t, err)
				}
				assert.Regexp(t, tt.want.expectedBody, string(body))
			}
		})
	}
}

func TestHandler_GetPing(t *testing.T) {
	// создадим конроллер моков и экземпляр мок-хранилища
	ctrl := gomock.NewController(t)
	s := mockStore.NewMockStorage(ctrl)

	// определим, какой результат будем получать от «хранилища»
	s.EXPECT().
		Ping(gomock.Any()).
		Return(nil)

	h := handlers.NewHandler(s, "http://localhost", "192.168.1.0/24")

	r := chi.NewRouter()
	r.Use(middlewares.Auth)
	r.Get("/ping", h.GetPing)

	srv := httptest.NewServer(r)
	defer srv.Close()

	type want struct {
		expectedCode int
	}

	tests := []struct {
		name   string
		method string
		want   want
	}{
		{
			name:   "method get success",
			method: http.MethodGet,
			want: want{
				expectedCode: http.StatusOK,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bodyReader io.Reader
			req, err := http.NewRequest(tt.method, srv.URL+"/ping", bodyReader)
			require.NoError(t, err)

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
