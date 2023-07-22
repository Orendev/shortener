package logger

import (
	"net/http"

	"go.uber.org/zap"
)

var Log *zap.Logger

type ResponseData struct {
	Status int
	Size   int
}

type LoggingResponseWriter struct {
	http.ResponseWriter
	ResponseData *ResponseData
}

// Write writes the data to the connection as part of an HTTP reply
func (r LoggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.ResponseData.Size += size
	return size, err
}

// WriteHeader sends an HTTP response header with the provided
func (r LoggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.ResponseData.Status = statusCode
}

// NewLogger конструктор создает глобавльную переменную Log
func NewLogger(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}

	cfg := zap.NewProductionConfig()

	cfg.Level = lvl

	zl, err := cfg.Build()
	if err != nil {
		return err
	}

	Log = zl
	return nil
}
