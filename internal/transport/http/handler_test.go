package http

import (
	"bytes"
	"github.com/Orendev/shortener/internal/config"
	"github.com/Orendev/shortener/internal/models"
	"github.com/Orendev/shortener/internal/random"
	service "github.com/Orendev/shortener/internal/services"
	"github.com/Orendev/shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlers_ShortLink(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}

	type fields struct {
		code string
		link string
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}

	tests := []struct {
		name    string
		fields  fields
		configs config.Configs
		args    args
		want    want
	}{
		{
			name: "positive test #1 handleShortLink",
			fields: fields{
				link: "https://google.com",
				code: random.Strn(8),
			},
			want: want{
				code:        http.StatusTemporaryRedirect,
				contentType: "text/plain",
			},
			configs: config.Configs{
				Host:            "",
				Port:            "8080",
				BaseURL:         "http://localhost:8080",
				Memory:          map[string]models.ShortLink{},
				FileStoragePath: "/tmp/test-short-url-db.json",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.configs.Memory = map[string]models.ShortLink{
				tt.fields.code: {
					ShortURL:    tt.fields.code,
					OriginalURL: tt.fields.link,
				},
			}

			file, err := storage.NewFile(&tt.configs)
			require.NoError(t, err)

			memoryStorage, err := storage.NewMemoryStorage(&tt.configs, nil, file)
			require.NoError(t, err)

			h := &Handler{
				shortLinkStorage: service.NewService(memoryStorage, &tt.configs),
			}

			defer func() {
				err = file.Remove()
				if err != nil {
					require.NoError(t, err)
				}
			}()

			req := httptest.NewRequest(http.MethodGet, "/"+tt.fields.code, nil)
			w := httptest.NewRecorder()

			h.ShortLink(w, req)

			res := w.Result()
			err = res.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, res.StatusCode, tt.want.code)
			assert.Equal(t, tt.fields.link, res.Header.Get("Location"))

		})
	}
}

func TestHandlers_ShortLinkAdd(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	type fields struct {
		url string
	}

	tests := []struct {
		name    string
		fields  fields
		configs config.Configs
		want    want
	}{
		{
			name: "positive test #1 handleShortLinkAdd",
			fields: fields{
				url: "https://google.com",
			},
			want: want{
				code:        http.StatusCreated,
				contentType: "text/plain",
				response:    "http://localhost:8080/",
			},
			configs: config.Configs{
				Host:            "",
				Port:            "8080",
				BaseURL:         "http://localhost:8080",
				Memory:          map[string]models.ShortLink{},
				FileStoragePath: "/tmp/test-short-url-db.json",
			},
		},
		{
			name: "negative test #2 handleShortLinkAdd",
			fields: fields{
				url: "",
			},
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				response:    "http://localhost:8080/",
			},
			configs: config.Configs{
				Host:            "",
				Port:            "8080",
				BaseURL:         "http://localhost:8080",
				Memory:          map[string]models.ShortLink{},
				FileStoragePath: "/tmp/test-short-url-db.json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.configs.Memory = map[string]models.ShortLink{}

			file, err := storage.NewFile(&tt.configs)
			require.NoError(t, err)

			memoryStorage, err := storage.NewMemoryStorage(&tt.configs, nil, file)
			require.NoError(t, err)

			h := &Handler{
				shortLinkStorage: service.NewService(memoryStorage, &tt.configs),
			}

			defer func() {
				err = file.Remove()
				if err != nil {
					require.NoError(t, err)
				}
			}()

			body := strings.NewReader(tt.fields.url)
			req := httptest.NewRequest(http.MethodPost, "/", body)
			w := httptest.NewRecorder()

			h.ShortLinkAdd(w, req)

			res := w.Result()

			assert.Equal(t, res.StatusCode, tt.want.code)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))

			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			err = res.Body.Close()
			require.NoError(t, err)

			for code, shortLink := range tt.configs.Memory {
				assert.Equal(t, tt.want.response+code, string(resBody))
				assert.Equal(t, tt.fields.url, shortLink.OriginalURL)
			}

		})
	}
}

func Test_handler_ApiShorten(t *testing.T) {

	cfg := config.Configs{
		Host:            "",
		Port:            "8080",
		BaseURL:         "http://localhost:8080",
		Memory:          map[string]models.ShortLink{},
		FileStoragePath: "/tmp/test-short-url-db.json",
	}

	memory := cfg.Memory

	file, err := storage.NewFile(&cfg)
	require.NoError(t, err)

	memoryStorage, err := storage.NewMemoryStorage(&cfg, nil, file)
	require.NoError(t, err)

	h := &Handler{
		shortLinkStorage: service.NewService(memoryStorage, &cfg),
	}

	handler := http.HandlerFunc(h.APIShorten)
	srv := httptest.NewServer(handler)

	defer srv.Close()
	defer func() {
		err = file.Remove()
		if err != nil {
			require.NoError(t, err)
		}
	}()

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
				expectedBody: "",
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

			req, err := http.NewRequest(tt.method, srv.URL, bodyReader)
			require.NoError(t, err)

			if len(tt.body) > 0 {
				req.Header.Set("Content-Type", "application/json")

			}

			resp, err := srv.Client().Do(req)
			require.NoError(t, err)

			defer func() {
				err := resp.Body.Close()
				if err != nil {
					require.NoError(t, err)
				}
			}()

			assert.NoError(t, err, "error making HTTP request")
			assert.Equal(t, tt.want.expectedCode, resp.StatusCode, "ShortLinkResponse code didn't match expected")

			// проверяем корректность полученного тела ответа, если мы его ожидаем
			if len(memory) > 0 && resp.StatusCode == http.StatusCreated {
				for _, link := range memory {
					tt.want.expectedBody = `{
						"result": "` + link.ShortURL + `"
					}`
					break
				}

				if tt.want.expectedBody != "" {
					body, err := io.ReadAll(resp.Body)
					if err != nil {
						require.NoError(t, err)
					}
					assert.JSONEq(t, tt.want.expectedBody, string(body))
				}
			}
		})
	}
}
