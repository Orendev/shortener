package compress

import (
	"compress/gzip"
	"io"
	"net/http"
)

// GzipWriter реализует интерфейс http.ResponseWriter и позволяет прозрачно для сервера
// сжимать передаваемые данные и выставлять правильные HTTP-заголовки
type GzipWriter struct {
	rw     http.ResponseWriter
	writer *gzip.Writer
}

// NewGzipWriter the constructor creates a GzipWriter
func NewGzipWriter(w http.ResponseWriter) *GzipWriter {
	return &GzipWriter{
		rw:     w,
		writer: gzip.NewWriter(w),
	}
}

// Header returns the header map that will be sent by
func (zw *GzipWriter) Header() http.Header {
	return zw.rw.Header()
}

// Write writes the data to the connection as part of an HTTP reply.
func (zw *GzipWriter) Write(b []byte) (int, error) {
	return zw.writer.Write(b)
}

// WriteHeader sends an HTTP response header with the provided
func (zw *GzipWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		zw.rw.Header().Set("Content-Encoding", "gzip")
	}
	zw.rw.WriteHeader(statusCode)
}

// Close закрывает gzip.Writer и досылает все данные из буфера.
func (zw *GzipWriter) Close() error {
	return zw.writer.Close()
}

// GzipReader реализует интерфейс io.ReadCloser и позволяет прозрачно для сервера
// декомпрессировать получаемые от клиента данные
type GzipReader struct {
	r      io.ReadCloser
	reader *gzip.Reader
}

// NewGzipReader декомпрессировать получаемые от клиента данные
func NewGzipReader(r io.ReadCloser) (*GzipReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &GzipReader{
		r:      r,
		reader: zr,
	}, nil
}

// Read reads up to len(p) bytes into p. It returns the number of bytes
func (zr GzipReader) Read(p []byte) (n int, err error) {
	return zr.reader.Read(p)
}

// Close is the interface that wraps the basic Close method.
func (zr GzipReader) Close() error {
	if err := zr.r.Close(); err != nil {
		return err
	}
	return zr.reader.Close()
}
