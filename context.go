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
}
