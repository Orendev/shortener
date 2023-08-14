package config

import (
	"encoding/json"
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_initFlag(t *testing.T) {
	type args struct {
		cfg      *Configs
		fs       *flag.FlagSet
		commands []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "commands successful",
			args: args{
				cfg: &Configs{},
				fs:  flag.NewFlagSet("shortener-test-initFlag", flag.ContinueOnError),
				commands: []string{
					"test",
					"-a", "Hello",
					"-b", "World",
					"-f", "/tmp/short-url-db.json",
					"-d", "host=localhost user=shortener password=secret dbname=shortener sslmode=disable",
					"-c", "./config/shortener.json",
					"-s", "true",
				},
			},
		},
		{
			name: "commands errors",
			args: args{
				cfg: &Configs{},
				fs:  flag.NewFlagSet("shortener-test-initFlag", flag.ContinueOnError),
				commands: []string{
					"test",
					"-tt", "Hello",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			os.Args = []string{}
			if len(tt.args.commands) > 0 {
				os.Args = tt.args.commands
			}
			if err := initFlag(tt.args.cfg, tt.args.fs); (err != nil) != tt.wantErr {
				t.Errorf("initFlag() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_initEnv(t *testing.T) {
	type args struct {
		cfg *Configs
		env map[string]string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "env successful",
			args: args{
				cfg: &Configs{},
				env: map[string]string{
					"ENABLE_HTTPS": "true",
				},
			},
		},
		{
			name: "env error",
			args: args{
				cfg: &Configs{},
				env: map[string]string{
					"ENABLE_HTTPS": "test",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.args.env) > 0 {
				for k, v := range tt.args.env {
					err := os.Setenv(k, v)
					require.NoError(t, err)
				}
			}

			if err := initEnv(tt.args.cfg); (err != nil) != tt.wantErr {
				t.Errorf("initEnv() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_setValueString(t *testing.T) {
	type args struct {
		value        string
		defaultValue string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "default value",
			args: args{
				value:        "",
				defaultValue: "default-test",
			},
			want: "default-test",
		},
		{
			name: "value",
			args: args{
				value:        "value-test",
				defaultValue: "default-test",
			},
			want: "value-test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := setValueString(tt.args.value, tt.args.defaultValue); got != tt.want {
				t.Errorf("setValueString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_initFile(t *testing.T) {

	type args struct {
		cfg      *Configs
		fs       *flag.FlagSet
		commands []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "config file successful",
			args: args{
				cfg: &Configs{},
				fs:  flag.NewFlagSet("shortener-test-initFile", flag.ContinueOnError),
				commands: []string{
					"test",
					"-a", "Hello",
					"-b", "World",
					"-f", "/tmp/short-url-db.json",
					"-c", "/tmp/shortener.json",
					"-d", "host=localhost user=shortener password=secret dbname=shortener sslmode=disable",
					"-s", "true",
				},
			},
		},
		{
			name: "config file error",
			args: args{
				cfg: &Configs{},
				fs:  flag.NewFlagSet("shortener-test-initFlag", flag.ContinueOnError),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = []string{}
			if len(tt.args.commands) > 0 {
				os.Args = tt.args.commands
				err := initFlag(tt.args.cfg, tt.args.fs)
				require.NoError(t, err)
			}

			if len(tt.args.cfg.Config) > 0 {
				file, err := os.OpenFile(tt.args.cfg.Config, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
				require.NoError(t, err)
				defer func() {
					if err := file.Close(); err != nil {
						require.NoError(t, err)
					}
				}()

				fileConfig := FileConfig{
					Addr:            "localhost:8080",
					BaseURL:         "localhost",
					DatabaseDSN:     "host=localhost user=shortener password=secret dbname=shortener sslmode=disable",
					IsHTTPS:         true,
					FileStoragePath: "/tmp/short-url-db.json",
				}

				writeData, err := json.Marshal(fileConfig)
				if err != nil {
					require.NoError(t, err)
				}
				_, err = file.Write(append(writeData, '\n'))
				if err != nil {
					require.NoError(t, err)
				}
			}

			if tt.wantErr {
				tt.args.cfg.Config = "test"
			}

			if err := initFile(tt.args.cfg, tt.args.fs); (err != nil) != tt.wantErr {
				t.Errorf("initFile() error = %v, wantErr %v", err, tt.wantErr)
			}

			if len(tt.args.cfg.Config) > 0 && !tt.wantErr {
				err := os.Remove(tt.args.cfg.Config)
				require.NoError(t, err)
			}

		})
	}
}

func Test_initDefaultValue(t *testing.T) {
	type args struct {
		cfg *Configs
	}
	type want struct {
		value        string
		defaultValue string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Default value successful",
			args: args{
				cfg: &Configs{},
			},
			want: want{
				defaultValue: "localhost:8080",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initDefaultValue(tt.args.cfg)
			assert.Equal(t, tt.want.defaultValue, tt.args.cfg.Server.Addr, "defaultValue didn't match expected")
		})
	}
}

func TestNew(t *testing.T) {

	cfg := Configs{
		Server: Server{
			Addr:    "Hello",
			IsHTTPS: true,
		},
		BaseURL: "World",
		File:    File{FileStoragePath: "/tmp/short-url-db.json"},
		Cert: Cert{
			CertFile: "cert.pem",
			KeyFile:  "key.pem",
		},
		Log: Log{FlagLogLevel: "info"},
		Database: Database{
			DatabaseDSN: "host=localhost user=shortener password=secret dbname=shortener sslmode=disable",
		},
		Config: "",
	}

	type args struct {
		cfg      *Configs
		commands []string
		env      map[string]string
	}

	tests := []struct {
		name string
		args args
		want *Configs
	}{
		{
			name: "commands successful",
			args: args{
				cfg: &Configs{},
				commands: []string{
					"test1",
					"-a", "Hello",
					"-b", "World",
					"-f", "/tmp/short-url-db.json",
					"-d", "host=localhost user=shortener password=secret dbname=shortener sslmode=disable",
					"-c", "",
					"-s", "true",
				},
				env: map[string]string{
					"ENABLE_HTTPS": "true",
				},
			},
			want: &cfg,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = []string{}
			if len(tt.args.commands) > 0 {
				os.Args = tt.args.commands
			}
			if len(tt.args.env) > 0 {
				for k, v := range tt.args.env {
					err := os.Setenv(k, v)
					require.NoError(t, err)
				}
			}

			got, err := New()
			require.NoError(t, err)
			assert.Equalf(t, tt.want, got, "New()")
		})
	}
}
