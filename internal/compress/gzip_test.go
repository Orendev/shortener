package compress_test

import (
	"bytes"
	"compress/gzip"
	"github.com/Orendev/shortener/internal/configs"
	"github.com/Orendev/shortener/internal/models"
	"github.com/Orendev/shortener/internal/routes"
	storage2 "github.com/Orendev/shortener/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var cfg = configs.Configs{
	Host:            "",
	Port:            "8080",
	BaseURL:         "http://localhost:8080",
	FileStoragePath: "/tmp/test-short-url-db.json",
}

func TestGzipMiddlewareSendsGzip(t *testing.T) {

	cfg.Memory = map[string]models.ShortLink{}

	memory := cfg.Memory

	fileDB, err := storage2.NewFileDB(&cfg)
	if err != nil {
		require.NoError(t, err)
	}

	memoryStorage, err := storage2.NewMemoryStorage(&cfg, fileDB)
	if err != nil {
		log.Fatal(err)
		return
	}

	r := routes.Routes(chi.NewRouter(), memoryStorage, &cfg)

	srv := httptest.NewServer(r)

	defer srv.Close()

	defer func() {
		err := fileDB.Remove()
		if err != nil {
			require.NoError(t, err)
		}
	}()

	type want struct {
		contentType     string
		contentEncoding string
		expectedCode    int
		expectedBody    string
	}

	tests := []struct {
		name   string // добавляем название тестов
		method string
		path   string // добавляем роут в табличные тесты
		body   string // добавляем тело запроса в табличные тесты
		want   want
	}{
		{
			name:   "method_post_success sends gzip",
			method: http.MethodPost,
			path:   "/api/shorten",
			body:   `{"url":"https://practicum.yandex.ru/"}`,
			want: want{
				contentEncoding: "gzip",
				expectedCode:    http.StatusCreated,
				expectedBody:    "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			zb := gzip.NewWriter(buf)
			_, err := zb.Write([]byte(tt.body))
			require.NoError(t, err)

			err = zb.Close()
			require.NoError(t, err)

			r := httptest.NewRequest(tt.method, srv.URL+tt.path, buf)
			r.RequestURI = ""
			if tt.want.contentEncoding != "" {
				r.Header.Set("Content-Encoding", tt.want.contentEncoding)
			}

			resp, err := http.DefaultClient.Do(r)
			require.NoError(t, err)
			require.Equal(t, tt.want.expectedCode, resp.StatusCode)

			defer func() {
				err := resp.Body.Close()
				if err != nil {
					require.NoError(t, err)
				}
			}()

			if len(memory) > 0 {
				for _, link := range memory {
					tt.want.expectedBody = `{
						"result": "` + link.Result + `"
					}`
					break
				}

				// проверяем корректность полученного тела ответа, если мы его ожидаем
				if tt.want.expectedBody != "" {
					b, err := io.ReadAll(resp.Body)
					require.NoError(t, err)
					require.JSONEq(t, tt.want.expectedBody, string(b))
				}
			}

		})
	}
}

func TestGzipMiddlewareAcceptsGzip(t *testing.T) {

	cfg.Memory = map[string]models.ShortLink{}

	memory := cfg.Memory

	fileDB, err := storage2.NewFileDB(&cfg)
	if err != nil {
		require.NoError(t, err)
	}

	memoryStorage, err := storage2.NewMemoryStorage(&cfg, fileDB)
	if err != nil {
		log.Fatal(err)
		return
	}

	r := routes.Routes(chi.NewRouter(), memoryStorage, &cfg)

	srv := httptest.NewServer(r)

	defer srv.Close()

	defer func() {
		err := fileDB.Remove()
		if err != nil {
			require.NoError(t, err)
		}
	}()

	type want struct {
		acceptEncoding string
		contentType    string
		expectedCode   int
		expectedBody   string
	}

	tests := []struct {
		name   string // добавляем название тестов
		method string
		path   string // добавляем роут в табличные тесты
		body   string // добавляем тело запроса в табличные тесты
		want   want
	}{
		{
			name:   "method_post_success accepts gzip",
			method: http.MethodPost,
			path:   "/api/shorten",
			body:   `{"url":"https://practicum.yandex.ru/"}`,
			want: want{
				acceptEncoding: "gzip",
				contentType:    "text/html",
				expectedCode:   http.StatusCreated,
				expectedBody:   "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.NewBufferString(tt.body)

			r := httptest.NewRequest(tt.method, srv.URL+tt.path, buf)
			r.RequestURI = ""
			if tt.want.acceptEncoding != "" {
				r.Header.Set("Accept-Encoding", tt.want.acceptEncoding)
			}

			if tt.want.contentType != "" {
				r.Header.Set("Content-Type", tt.want.contentType)
			}

			resp, err := http.DefaultClient.Do(r)
			require.NoError(t, err)
			require.Equal(t, tt.want.expectedCode, resp.StatusCode)

			defer func() {
				err := resp.Body.Close()
				if err != nil {
					require.NoError(t, err)
				}
			}()

			zr, err := gzip.NewReader(resp.Body)
			require.NoError(t, err)

			b, err := io.ReadAll(zr)
			require.NoError(t, err)

			if len(memory) > 0 {
				for _, link := range memory {
					tt.want.expectedBody = `{
						"result": "` + link.Result + `"
					}`
					break
				}
			}

			// проверяем корректность полученного тела ответа, если мы его ожидаем
			if tt.want.expectedBody != "" {
				require.JSONEq(t, tt.want.expectedBody, string(b))
			}

		})
	}

}
