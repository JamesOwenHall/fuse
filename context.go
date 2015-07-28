package fuse

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type Context struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	Params         map[string]string
	Form           url.Values
	PostForm       url.Values
	Data           map[string]interface{}

	engine       *Engine
	handlerIndex int
	handler      Handler
}

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

func (c *Context) Text(code int, text string) {
	c.ResponseWriter.Header().Set("Content-Type", "text/plain")
	c.ResponseWriter.WriteHeader(code)
	c.ResponseWriter.Write([]byte(text))
}

func (c *Context) TextOk(text string) {
	c.Text(http.StatusOK, text)
}

func (c *Context) Html(code int, name string) {
	c.ResponseWriter.WriteHeader(code)
	c.engine.execTemplate(c, name)
}

func (c *Context) HtmlOk(name string) {
	c.Html(http.StatusOK, name)
}

func (c *Context) Json(code int) {
	c.ResponseWriter.Header().Set("Content-Type", "application/json")
	c.ResponseWriter.WriteHeader(code)
	encoder := json.NewEncoder(c.ResponseWriter)
	encoder.Encode(c.Data)
}

func (c *Context) JsonOk() {
	c.Json(http.StatusOK)
}

func (c *Context) SeeOther(location string) {
	c.ResponseWriter.Header().Set("Location", location)
	c.ResponseWriter.WriteHeader(http.StatusSeeOther)
}
