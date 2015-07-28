package fuse

import (
	"encoding/json"
	"net/http"
	"net/url"
)

// Context holds all necessary information about an incoming request.
type Context struct {
	// Request is the original http.Request.
	Request *http.Request
	// ResponseWriter is the original http.ResponseWriter.
	ResponseWriter http.ResponseWriter
	// Params is a map of the path parameters.
	Params map[string]string
	// Form holds all of the form data from the URL query and the POST/PUT
	// data.
	Form url.Values
	// PostForm holds only the POST/PUT data.
	PostForm url.Values
	// InData provides a way for middleware to pass data to handlers.
	InData map[string]interface{}
	// OutData is used to render dynamic content.  Calls to Context.Html and
	// Context.Json use this data.
	OutData map[string]interface{}
	// Session contains the current request's session data.
	Session map[interface{}]interface{}

	engine       *Engine
	handlerIndex int
	handler      Handler
}

// Next calls the next middleware in the chain.  If called from the handler,
// it will panic.
func (c *Context) Next() {
	handlerLen := len(c.engine.middleware)
	if c.handlerIndex == handlerLen {
		c.handler(c)
	} else if c.handlerIndex < handlerLen {
		currentIndex := c.handlerIndex
		c.handlerIndex++
		c.engine.middleware[currentIndex](c)
	} else {
		panic("Can't call Context.Next() on handler")
	}
}

// Text writes the text to the response with the given HTTP code.
func (c *Context) Text(code int, text string) {
	c.ResponseWriter.Header().Set("Content-Type", "text/plain")
	c.ResponseWriter.WriteHeader(code)
	c.ResponseWriter.Write([]byte(text))
}

// TextOk writes the text to the response with the 200 OK HTTP code.
func (c *Context) TextOk(text string) {
	c.Text(http.StatusOK, text)
}

// Html executes the named template with the given HTTP code.  It uses the
// context's OutData in the template.
func (c *Context) Html(code int, name string) {
	c.ResponseWriter.WriteHeader(code)
	c.engine.execTemplate(c, name)
}

// Html executes the named template with the 200 OK HTTP code.  It uses the
// context's OutData in the template.
func (c *Context) HtmlOk(name string) {
	c.Html(http.StatusOK, name)
}

// Json encodes the context's OutData to JSON with the given HTTP code.
func (c *Context) Json(code int) {
	c.ResponseWriter.Header().Set("Content-Type", "application/json")
	c.ResponseWriter.WriteHeader(code)
	encoder := json.NewEncoder(c.ResponseWriter)
	encoder.Encode(c.OutData)
}

// Json encodes the context's OutData to JSON with the 200 OK HTTP code.
func (c *Context) JsonOk() {
	c.Json(http.StatusOK)
}

// SeeOther redirects the request to the location using the 303 See Other HTTP
// code.
func (c *Context) SeeOther(location string) {
	c.ResponseWriter.Header().Set("Location", location)
	c.ResponseWriter.WriteHeader(http.StatusSeeOther)
}
