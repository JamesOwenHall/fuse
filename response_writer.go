package fuse

import (
	"net/http"

	"github.com/gorilla/sessions"
)

// responseWriter is a wrapper for an http.ResponseWriter that saves an
// Engine's session before writing.
type responseWriter struct {
	writer  http.ResponseWriter
	request *http.Request
	session *sessions.Session
	saved   bool
}

func (r *responseWriter) Header() http.Header {
	return r.writer.Header()
}

func (r *responseWriter) Write(b []byte) (int, error) {
	if !r.saved {
		r.session.Save(r.request, r.writer)
		r.saved = true
	}

	return r.writer.Write(b)
}

func (r *responseWriter) WriteHeader(code int) {
	if !r.saved {
		r.session.Save(r.request, r.writer)
		r.saved = true
	}

	r.writer.WriteHeader(code)
}
