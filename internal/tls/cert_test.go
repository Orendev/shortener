package tls

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_doesFileExist(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "file not found",
			args: args{
				path: "/tmp/test.pem",
			},
			want: false,
		},
		{
			name: "the file exists",
			args: args{
				path: "/tmp/test.pem",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got bool
			if tt.want {
				err := os.WriteFile(tt.args.path, []byte(""), 664)
				require.NoError(t, err)
				got = doesFileExist(tt.args.path)
				err = os.Remove(tt.args.path)
				require.NoError(t, err)
			} else {
				got = doesFileExist(tt.args.path)
			}

			if got != tt.want {
				t.Errorf("doesFileExist() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generate(t *testing.T) {
	type args struct {
		certFile string
		keyFile  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "successful certificate generation",
			args: args{
				certFile: "/tmp/cert.pem",
				keyFile:  "/tmp/key.pem",
			},
			wantErr: false,
		},
		{
			name: "error generating certificates",
			args: args{
				certFile: "",
				keyFile:  "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := generate(tt.args.certFile, tt.args.keyFile); (err != nil) != tt.wantErr {
				t.Errorf("generate() error = %v, wantErr %v", err, tt.wantErr)
			}

			if len(tt.args.certFile) > 0 {
				err := os.Remove(tt.args.certFile)
				require.NoError(t, err)
			}

			if len(tt.args.keyFile) > 0 {
				err := os.Remove(tt.args.keyFile)
				require.NoError(t, err)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		certFile string
		keyFile  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "successful initialization of certificates",
			args: args{
				certFile: "/tmp/cert.pem",
				keyFile:  "/tmp/key.pem",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := New(tt.args.certFile, tt.args.keyFile); (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
			}

			if len(tt.args.certFile) > 0 {
				err := os.Remove(tt.args.certFile)
				require.NoError(t, err)
			}

			if len(tt.args.keyFile) > 0 {
				err := os.Remove(tt.args.keyFile)
				require.NoError(t, err)
			}
		})
	}
}
