package server

import (
	"github.com/Orendev/shortener/internal/api"
	"github.com/Orendev/shortener/internal/app/repository/shortlinks"
	"github.com/Orendev/shortener/internal/pkg/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
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
		name   string
		fields fields
		args   args
		want   want
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
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shortLinkStore, _ := shortlinks.NewStorage(map[string]shortlinks.ShortLink{
				tt.fields.code: {
					Code: tt.fields.code,
					Link: tt.fields.link,
				},
			})
			sl, _ := shortlinks.New(*shortLinkStore)
			a, _ := api.New(sl)
			h := &Handlers{
				api: a,
			}

			req := httptest.NewRequest(http.MethodGet, "/"+tt.fields.code, nil)
			w := httptest.NewRecorder()

			h.handleShortLink(w, req)

			res := w.Result()

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
		name   string
		fields fields
		want   want
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := map[string]shortlinks.ShortLink{}
			shortLinkStore, _ := shortlinks.NewStorage(data)

			sl, _ := shortlinks.New(*shortLinkStore)
			a, _ := api.New(sl)
			h := &Handlers{
				api: a,
			}
			body := strings.NewReader(tt.fields.link)
			req := httptest.NewRequest(http.MethodPost, "/", body)
			w := httptest.NewRecorder()

			h.handleShortLinkAdd(w, req)

			res := w.Result()

			assert.Equal(t, res.StatusCode, tt.want.code)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))

			resBody, err := ioutil.ReadAll(res.Body)
			require.NoError(t, err)
			err = res.Body.Close()
			require.NoError(t, err)

			for code, shortLink := range data {
				assert.Equal(t, tt.want.response+code, string(resBody))
				assert.Equal(t, tt.fields.link, shortLink.Link)
			}

		})
	}
}
