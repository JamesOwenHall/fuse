package fuse

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Handler func(c *Context)

type Engine struct {
	router     *httprouter.Router
	middleware []Handler
}

func New() *Engine {
	return &Engine{
		router:     httprouter.New(),
		middleware: make([]Handler, 0),
	}
}

func (e *Engine) Run(addr string) {
	http.ListenAndServe(addr, e.router)
}

func (e *Engine) Use(handler Handler) {
	e.middleware = append(e.middleware, handler)
}

func (e *Engine) GET(path string, handler Handler) {
	e.router.GET(path, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		r.ParseForm()

		params := make(map[string]string)
		for _, param := range p {
			params[param.Key] = param.Value
		}

		c := Context{
			Request:        r,
			ResponseWriter: responseWriter{w},
			Params:         params,
			Form:           r.Form,
			PostForm:       r.PostForm,
			engine:         e,
			handler:        handler,
		}

		c.Next()
	})
}
