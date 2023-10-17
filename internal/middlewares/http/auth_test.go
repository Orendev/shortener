package http

import (
	"context"
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
			name:   "positive test auth middleware",
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

func Test_generateAuthHeaderFromToken(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "getting a successful token header",
			args: args{
				token: "test-token",
			},
			want: "Bearer test-token",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, generateAuthHeaderFromToken(tt.args.token), "generateAuthHeaderFromToken(%v)", tt.args.token)
		})
	}
}

func Test_extractTokenFromAuthHeader(t *testing.T) {
	type args struct {
		val string
	}
	tests := []struct {
		name      string
		args      args
		wantToken string
		wantOk    bool
	}{
		{
			name: "successful extraction of token from from authheader",
			args: args{
				val: "bearer test-token",
			},
			wantToken: "test-token",
			wantOk:    true,
		},

		{
			name: "failed to extract token from from authheader",
			args: args{
				val: "test-token",
			},
			wantToken: "",
			wantOk:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToken, gotOk := extractTokenFromAuthHeader(tt.args.val)
			assert.Equalf(t, tt.wantToken, gotToken, "extractTokenFromAuthHeader(%v)", tt.args.val)
			assert.Equalf(t, tt.wantOk, gotOk, "extractTokenFromAuthHeader(%v)", tt.args.val)
		})
	}
}

func TestHTTPToContext(t *testing.T) {

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	require.NoError(t, err)
	ctx, err := NewSigner(context.Background())
	require.NoError(t, err)

	tokenString, ok := ctx.Value(auth.JwtContextKey).(string)
	if ok {
		req.Header.Set(auth.HeaderAuthorizationKey, "bearer "+tokenString)
	}

	type args struct {
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "HTTP to Context",
			args: args{
				r: req,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := HTTPToContext(tt.args.r)
			require.NoError(t, err)

		})
	}
}

func TestNewSigner(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		{
			name: "creates a JWT by specifying the key ID",
			args: args{
				context.Background(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSigner(tt.args.ctx)
			require.NoError(t, err)
			_, ok := got.Value(auth.JwtContextKey).(string)
			assert.Equal(t, ok, true)
			assert.NotEqualf(t, tt.want, got, "NewSigner(%v)", tt.args.ctx)
		})
	}
}
