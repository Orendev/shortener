package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGzip(t *testing.T) {

	type want struct {
		statusCode          int
		contentType         string
		expectedContentType string
		headerContent       string
		expectedBody        string
		acceptEncoding      string
		contentEncoding     string
	}
	tests := []struct {
		name   string
		method string
		body   string // добавляем тело запроса в табличные тесты
		want   want
	}{
		{
			name:   "Positive test#1",
			method: http.MethodPost,
			body:   "{\"url\": \"ya.ru\"}",
			want: want{
				statusCode:          http.StatusOK,
				contentType:         "application/json",
				expectedContentType: "application/x-gzip",
				headerContent:       "gzip",
				expectedBody:        "",
				acceptEncoding:      "gzip",
				contentEncoding:     "",
			},
		},
		{
			name:   "Positive test#2",
			method: http.MethodPost,
			want: want{
				statusCode:          http.StatusOK,
				contentType:         "application/json",
				expectedContentType: "text/plain; charset=utf-8",
				headerContent:       "",
				expectedBody:        "",
				acceptEncoding:      "",
				contentEncoding:     "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(tt.method, "/", nil)

			req.Header.Set("Accept-Encoding", tt.want.acceptEncoding)
			req.Header.Set("Content-Type", tt.want.contentType)
			req.Header.Set("Content-Encoding", tt.want.contentEncoding)
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, err := w.Write([]byte(tt.body))
				if err != nil {
					require.NoError(t, err)
				}
			})

			Gzip(handler).ServeHTTP(rr, req)

			if rr.Header().Get("Content-Encoding") != tt.want.contentEncoding {
				t.Errorf("expected response %s, but got Content-Encoding header '%s'", tt.want.contentEncoding, rr.Header().Get("Content-Encoding"))
			}
			if rr.Header().Get("Content-Type") != tt.want.expectedContentType {
				t.Errorf("expected response %s, but got Content-Type header %s", tt.want.expectedContentType, rr.Header().Get("Content-Type"))
			}
			if rr.Code != tt.want.statusCode {
				t.Errorf("expected response code %d, but got %d", tt.want.statusCode, rr.Code)
			}
		})
	}
}
