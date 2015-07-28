package fuse

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
)

// Handler is the type of function used to handle all requests in Fuse,
// including middleware.
type Handler func(c *Context)

// Engine is the Fuse server.  It should always be created using the fuse.New
// function.
type Engine struct {
	// NotFound is called whenever a request comes in that isn't mapped to a
	// handler and doesn't correspond to a file name in the public directory.
	NotFound Handler
	// PublicDir is the directory which holds all files that are publicly
	// accessible.  The default is "public".
	PublicDir string
	// TempalteGlob is the glob that defines which files to parse as HTML
	// templates.  The default is "templates/*.tpl"
	TemplateGlob string

	router       *httprouter.Router
	middleware   []Handler
	sessionStore sessions.Store
}

// New returns an initialized instance of *Engine.
func New(sessionSecret []byte) *Engine {
	engine := &Engine{
		PublicDir:    "public",
		TemplateGlob: "templates/*.tpl",
		middleware:   make([]Handler, 0),
		sessionStore: sessions.NewCookieStore(sessionSecret),
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

// Run calls ListenAndServe.
func (e *Engine) Run(addr string) {
	http.ListenAndServe(addr, context.ClearHandler(e.router))
}

// Use adds the handler to the end of the middleware chain.
func (e *Engine) Use(handler Handler) {
	e.middleware = append(e.middleware, handler)
}

// GET defines a route for GET requests to the handler.  The path can contain
// parameters prefixed with a colon (:).  For example, the following requests
// would all be routed to the handler if we define the path as /user/:name
//     /user/foo
//     /user/bar
//     /user/baz
// However, the following paths would not be routed.
//     /user
//     /user/
//     /user/foo/profile
//     /profile/user/foo
func (e *Engine) GET(path string, handler Handler) {
	e.router.GET(path, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		c := e.makeContext(w, r, p, handler)
		c.Next()
	})
}

// POST defines a route for POST requests to the handler.  See GET for
// information about path parameters.
func (e *Engine) POST(path string, handler Handler) {
	e.router.POST(path, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		c := e.makeContext(w, r, p, handler)
		c.Next()
	})
}

// PUT defines a route for PUT requests to the handler.  See GET for
// information about path parameters.
func (e *Engine) PUT(path string, handler Handler) {
	e.router.PUT(path, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		c := e.makeContext(w, r, p, handler)
		c.Next()
	})
}

// DELETE defines a route for DELETE requests to the handler.  See GET for
// information about path parameters.
func (e *Engine) DELETE(path string, handler Handler) {
	e.router.DELETE(path, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		c := e.makeContext(w, r, p, handler)
		c.Next()
	})
}

// HEAD defines a route for HEAD requests to the handler.  See GET for
// information about path parameters.
func (e *Engine) HEAD(path string, handler Handler) {
	e.router.HEAD(path, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		c := e.makeContext(w, r, p, handler)
		c.Next()
	})
}

func (e *Engine) execTemplate(c *Context, name string) {
	tmpl, err := template.ParseGlob(e.TemplateGlob)
	if err != nil {
		c.Text(http.StatusInternalServerError, "Error parsing template: "+err.Error())
		return
	}

	tmpl.ExecuteTemplate(c.ResponseWriter, name, c.OutData)
}

func (e *Engine) makeContext(w http.ResponseWriter, r *http.Request, p httprouter.Params, handler Handler) *Context {
	r.ParseForm()

	params := make(map[string]string)
	for _, param := range p {
		params[param.Key] = param.Value
	}

	session, err := e.sessionStore.Get(r, "default")
	if err != nil {
		panic(err)
	}

	return &Context{
		Request:        r,
		ResponseWriter: &responseWriter{w, r, session, false},
		Params:         params,
		Form:           r.Form,
		PostForm:       r.PostForm,
		InData:         make(map[string]interface{}),
		OutData:        make(map[string]interface{}),
		Session:        session.Values,
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
