package fuse

import (
	"net/http"
)

type responseWriter struct {
	writer http.ResponseWriter
}

func (r responseWriter) Header() http.Header {
	return r.writer.Header()
}

func (r responseWriter) Write(b []byte) (int, error) {
	return r.writer.Write(b)
}

func (r responseWriter) WriteHeader(code int) {
	r.writer.WriteHeader(code)
}
