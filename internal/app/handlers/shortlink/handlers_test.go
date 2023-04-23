package shortlink

import (
	"bytes"
	models "github.com/Orendev/shortener/internal/app/models/shortlink"
	"github.com/Orendev/shortener/internal/app/repository/shortlink"
	service "github.com/Orendev/shortener/internal/app/service/shortlink"
	"github.com/Orendev/shortener/internal/configs"
	"github.com/Orendev/shortener/internal/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlers_handleShortLink(t *testing.T) {
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
		configs configs.Configs
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
			configs: configs.Configs{
				Host:    "",
				Port:    "8080",
				BaseURL: "http://localhost:8080",
				Memory:  map[string]models.ShortLink{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.configs.Memory = map[string]models.ShortLink{
				tt.fields.code: {
					Code: tt.fields.code,
					Link: tt.fields.link,
				},
			}
			memoryStorage, _ := shortlink.NewMemoryStorage(&tt.configs)

			h := &handler{
				ShortLinkRepository: service.NewService(memoryStorage, &tt.configs),
			}

			req := httptest.NewRequest(http.MethodGet, "/"+tt.fields.code, nil)
			w := httptest.NewRecorder()

			h.handleShortLink(w, req)

			res := w.Result()
			err := res.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, res.StatusCode, tt.want.code)
			assert.Equal(t, tt.fields.link, res.Header.Get("Location"))
		})
	}
}

func TestHandlers_handleShortLinkAdd(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	type fields struct {
		link string
	}

	tests := []struct {
		name    string
		fields  fields
		configs configs.Configs
		want    want
	}{
		{
			name: "positive test #1 handleShortLinkAdd",
			fields: fields{
				link: "https://google.com",
			},
			want: want{
				code:        http.StatusCreated,
				contentType: "text/plain",
				response:    "http://localhost:8080/",
			},
			configs: configs.Configs{
				Host:    "",
				Port:    "8080",
				BaseURL: "http://localhost:8080",
				Memory:  map[string]models.ShortLink{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.configs.Memory = map[string]models.ShortLink{}
			memoryStorage, _ := shortlink.NewMemoryStorage(&tt.configs)

			h := &handler{
				ShortLinkRepository: service.NewService(memoryStorage, &tt.configs),
			}
			body := strings.NewReader(tt.fields.link)
			req := httptest.NewRequest(http.MethodPost, "/", body)
			w := httptest.NewRecorder()

			h.handleShortLinkAdd(w, req)

			res := w.Result()

			assert.Equal(t, res.StatusCode, tt.want.code)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))

			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			err = res.Body.Close()
			require.NoError(t, err)

			for code, shortLink := range tt.configs.Memory {
				assert.Equal(t, tt.want.response+code, string(resBody))
				assert.Equal(t, tt.fields.link, shortLink.Link)
			}

		})
	}
}

func Test_handler_handleApiShorten(t *testing.T) {

	cfg := configs.Configs{
		Host:    "",
		Port:    "8080",
		BaseURL: "http://localhost:8080",
		Memory:  map[string]models.ShortLink{},
	}

	memoryStorage, _ := shortlink.NewMemoryStorage(&cfg)

	memory := cfg.Memory

	h := &handler{
		ShortLinkRepository: service.NewService(memoryStorage, &cfg),
	}

	handler := http.HandlerFunc(h.handleAPIShorten)
	srv := httptest.NewServer(handler)

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

			if len(tt.body) > 0 {
				req.Header.Set("Content-Type", "application/json")

			}

			resp, err := srv.Client().Do(req)
			require.NoError(t, err)

			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {

				}
			}(resp.Body)

			assert.NoError(t, err, "error making HTTP request")
			assert.Equal(t, tt.want.expectedCode, resp.StatusCode, "Response code didn't match expected")

			// проверяем корректность полученного тела ответа, если мы его ожидаем
			if len(memory) > 0 {
				for _, link := range memory {
					tt.want.expectedBody = `{
						"result": "` + link.Result + `"
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
