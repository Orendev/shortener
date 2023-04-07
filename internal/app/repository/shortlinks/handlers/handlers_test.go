package handlers

import (
	"github.com/Orendev/shortener/internal/app/repository/shortlinks"
	"github.com/Orendev/shortener/internal/app/repository/shortlinks/Model"
	"github.com/Orendev/shortener/internal/app/repository/shortlinks/storage"
	"github.com/Orendev/shortener/internal/configs"
	"github.com/Orendev/shortener/internal/pkg/random"
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
				Memory:  map[string]Model.ShortLink{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.configs.Memory = map[string]Model.ShortLink{
				tt.fields.code: {
					Code: tt.fields.code,
					Link: tt.fields.link,
				},
			}
			shortLinkStore, _ := storage.New(&tt.configs)
			sl, _ := shortlinks.New(shortLinkStore, &tt.configs)

			h := &handler{
				ShortLinkService: sl,
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
				Memory:  map[string]Model.ShortLink{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.configs.Memory = map[string]Model.ShortLink{}
			shortLinkStore, _ := storage.New(&tt.configs)

			sl, _ := shortlinks.New(shortLinkStore, &tt.configs)

			h := &handler{
				ShortLinkService: sl,
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
