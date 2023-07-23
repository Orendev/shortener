package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Orendev/shortener/internal/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuth(t *testing.T) {
	type args struct {
		next http.Handler
	}
	type want struct {
		expectedCode int
	}

	tests := []struct {
		name   string
		method string
		args   args
		want   want
	}{
		{
			name:   "positive test logger middleware",
			method: http.MethodGet,
			args: args{
				next: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}),
			},
			want: want{
				expectedCode: http.StatusOK,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ts := httptest.NewServer(Auth(tt.args.next))
			defer ts.Close()

			req, err := http.NewRequest(tt.method, ts.URL, nil)
			req.Header.Set(auth.HeaderAuthorizationKey, "fff1")
			if err != nil {
				require.NoError(t, err)
			}

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				require.NoError(t, err)
			}

			defer func() {
				err := resp.Body.Close()
				if err != nil {
					require.NoError(t, err)
				}
			}()

			assert.Regexp(t, "Bearer *", resp.Header.Get(auth.HeaderAuthorizationKey), "error http authorization header type for JWT")
			assert.Equal(t, tt.method, resp.Request.Method, "method didn't match expected")
			assert.Equal(t, tt.want.expectedCode, resp.StatusCode, "code didn't match expected")

		})
	}
}
