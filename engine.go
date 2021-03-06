package fuse

import (
	"html/template"
	"net/http"

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
	// Panic is called whenever a request panics.
	Panic Handler
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
		NotFound:              &fileServer{engine},
	}

	engine.NotFound = func(c *Context) {
		http.NotFound(c.ResponseWriter, c.Request)
	}

	engine.Panic = func(c *Context) {
		c.Text(http.StatusInternalServerError, "Internal server error")
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
		defer e.recoverFromPanic(c)
		c.Next()
	})
}

// POST defines a route for POST requests to the handler.  See GET for
// information about path parameters.
func (e *Engine) POST(path string, handler Handler) {
	e.router.POST(path, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		c := e.makeContext(w, r, p, handler)
		defer e.recoverFromPanic(c)
		c.Next()
	})
}

// PUT defines a route for PUT requests to the handler.  See GET for
// information about path parameters.
func (e *Engine) PUT(path string, handler Handler) {
	e.router.PUT(path, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		c := e.makeContext(w, r, p, handler)
		defer e.recoverFromPanic(c)
		c.Next()
	})
}

// DELETE defines a route for DELETE requests to the handler.  See GET for
// information about path parameters.
func (e *Engine) DELETE(path string, handler Handler) {
	e.router.DELETE(path, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		c := e.makeContext(w, r, p, handler)
		defer e.recoverFromPanic(c)
		c.Next()
	})
}

// HEAD defines a route for HEAD requests to the handler.  See GET for
// information about path parameters.
func (e *Engine) HEAD(path string, handler Handler) {
	e.router.HEAD(path, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		c := e.makeContext(w, r, p, handler)
		defer e.recoverFromPanic(c)
		c.Next()
	})
}

func (e *Engine) execTemplate(c *Context, name string) {
	tmpl, err := template.ParseGlob(e.TemplateGlob)
	if err != nil {
		panic(err)
	}

	err = tmpl.ExecuteTemplate(c.ResponseWriter, name, c.OutData)
	if err != nil {
		panic(err)
	}
}

func (e *Engine) makeContext(w http.ResponseWriter, r *http.Request, p httprouter.Params, handler Handler) *Context {
	r.ParseForm()

	params := make(map[string]string)
	for _, param := range p {
		params[param.Key] = param.Value
	}

	session, _ := e.sessionStore.Get(r, "default")

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

func (e *Engine) recoverFromPanic(c *Context) {
	if r := recover(); r != nil {
		e.Panic(c)
	}
}
