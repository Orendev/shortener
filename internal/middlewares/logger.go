package middlewares

import (
	"net/http"
	"time"

	"github.com/Orendev/shortener/internal/logger"
	"go.uber.org/zap"
)

// responseData структура, хранения сведения об ответе.
type responseData struct {
	status int
	size   int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

// Write writes the data to the connection as part of an HTTP reply.
func (r loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

// WriteHeader sends an HTTP response header with the provided.
func (r loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// Logger  middleware to log requests to the server.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}

		// создаем свой собственный ResponseWriter.
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		// внедряем оригинальную реализацию http.ResponseWriter.
		next.ServeHTTP(&lw, r)

		duration := time.Since(start)

		logger.Log.Info("got incoming HTTP request and response",
			zap.String("uri", r.RequestURI),
			zap.String("contentEncoding", r.Header.Get("Content-Encoding")),
			zap.Any("acceptEncoding", r.Header.Values("Accept-Encoding")),
			zap.String("method", r.Method),
			zap.Int("status", responseData.status),
			zap.Duration("duration", duration),
			zap.Int("size", responseData.size),
		)

	})
}
