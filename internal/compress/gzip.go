package compress

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// gzipWriter реализует интерфейс http.ResponseWriter и позволяет прозрачно для сервера
// сжимать передаваемые данные и выставлять правильные HTTP-заголовки
type gzipWriter struct {
	rw     http.ResponseWriter
	writer *gzip.Writer
}

func newGzipWriter(w http.ResponseWriter) *gzipWriter {
	return &gzipWriter{
		rw:     w,
		writer: gzip.NewWriter(w),
	}
}

func (zw *gzipWriter) Header() http.Header {
	return zw.rw.Header()
}

func (zw *gzipWriter) Write(b []byte) (int, error) {
	return zw.writer.Write(b)
}

func (zw *gzipWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		zw.rw.Header().Set("Content-Encoding", "gzip")
	}
	zw.rw.WriteHeader(statusCode)
}

// Close закрывает gzip.Writer и досылает все данные из буфера.
func (zw *gzipWriter) Close() error {
	return zw.writer.Close()
}

// gzipReader реализует интерфейс io.ReadCloser и позволяет прозрачно для сервера
// декомпрессировать получаемые от клиента данные
type gzipReader struct {
	r      io.ReadCloser
	reader *gzip.Reader
}

func newGzipReader(r io.ReadCloser) (*gzipReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &gzipReader{
		r:      r,
		reader: zr,
	}, nil
}

func (zr gzipReader) Read(p []byte) (n int, err error) {
	return zr.reader.Read(p)
}

func (zr *gzipReader) Close() error {
	if err := zr.r.Close(); err != nil {
		return err
	}
	return zr.reader.Close()
}

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// по умолчанию устанавливаем оригинальный http.ResponseWriter как тот,
		// который будем передавать следующей функции
		ow := w

		// проверяем, что клиент умеет получать от сервера сжатые данные в формате gzip
		acceptEncoding := r.Header.Get("Accept-Encoding")
		contentType := r.Header.Get("Content-Type")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")

		if supportsGzip && (contentType == "text/html" || contentType == "application/json") {
			// оборачиваем оригинальный http.ResponseWriter новым с поддержкой сжатия
			cw := newGzipWriter(w)
			// меняем оригинальный http.ResponseWriter на новый
			ow = cw
			// не забываем отправить клиенту все сжатые данные после завершения middleware
			defer func() {
				err := cw.Close()
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
				}
			}()
		}

		// проверяем, что клиент отправил серверу сжатые данные в формате gzip
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			// оборачиваем тело запроса в io.Reader с поддержкой декомпрессии
			cr, err := newGzipReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			// меняем тело запроса на новое
			r.Body = cr

			defer func() {
				err := cr.Close()
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
				}
			}()
		}

		next.ServeHTTP(ow, r)
	})
}
