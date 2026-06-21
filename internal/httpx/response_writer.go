package httpx

import (
	"net/http"
)

type responseWriter struct {
	http.ResponseWriter

	Status          int
	Bytes           int
	ResponseStarted bool
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
	}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.Status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.ResponseStarted = true
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.Status == 0 {
		rw.Status = http.StatusOK
	}

	n, err := rw.ResponseWriter.Write(b)

	rw.Bytes += n

	return n, err
}
