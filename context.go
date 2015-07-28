package fuse

import (
	"net/http"
	"net/url"
)

type Context struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	Params         map[string]string
	Form           url.Values
	PostForm       url.Values

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
	c.ResponseWriter.WriteHeader(code)
	c.ResponseWriter.Write([]byte(text))
}

func (c *Context) TextOk(text string) {
	c.Text(http.StatusOK, text)
}
