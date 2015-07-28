package fuse

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
)

type Handler func(c *Context)

type Engine struct {
	NotFound  Handler
	PublicDir string

	router     *httprouter.Router
	middleware []Handler
}

func New() *Engine {
	engine := &Engine{
		PublicDir:  "public",
		middleware: make([]Handler, 0),
	}

	engine.router = &httprouter.Router{
		RedirectTrailingSlash: true,
		NotFound:              &notFound{engine},
	}

	engine.NotFound = func(c *Context) {
		http.NotFound(c.ResponseWriter, c.Request)
	}

	return engine
}

func (e *Engine) Run(addr string) {
	http.ListenAndServe(addr, e.router)
}

func (e *Engine) Use(handler Handler) {
	e.middleware = append(e.middleware, handler)
}

func (e *Engine) GET(path string, handler Handler) {
	e.router.GET(path, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		c := e.makeContext(w, r, p, handler)
		c.Next()
	})
}

func (e *Engine) makeContext(w http.ResponseWriter, r *http.Request, p httprouter.Params, handler Handler) *Context {
	r.ParseForm()

	params := make(map[string]string)
	for _, param := range p {
		params[param.Key] = param.Value
	}

	return &Context{
		Request:        r,
		ResponseWriter: responseWriter{w},
		Params:         params,
		Form:           r.Form,
		PostForm:       r.PostForm,
		engine:         e,
		handler:        handler,
	}
}

func (e *Engine) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	filename := filepath.Join(e.PublicDir, r.URL.Path)
	if _, err := os.Stat(filename); err == nil {
		c := e.makeContext(w, r, nil, func(c *Context) {
			http.ServeFile(c.ResponseWriter, c.Request, filename)
		})
		c.Next()
	} else {
		c := e.makeContext(w, r, nil, e.NotFound)
		c.Next()
	}
}
