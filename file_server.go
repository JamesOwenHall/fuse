package fuse

import (
	"net/http"
	"os"
	"path/filepath"
)

// fileServer checks if a file exists in an Engine's public directory and
// serves it.  Otherwise, it call the Engine's NotFound handler.
type fileServer struct {
	e *Engine
}

func (f *fileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	filename := filepath.Join(f.e.PublicDir, r.URL.Path)
	if _, err := os.Stat(filename); err == nil {
		c := f.e.makeContext(w, r, nil, func(c *Context) {
			http.ServeFile(c.ResponseWriter, c.Request, filename)
		})
		c.Next()
	} else {
		c := f.e.makeContext(w, r, nil, f.e.NotFound)
		c.Next()
	}
}
