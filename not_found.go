package fuse

import (
	"net/http"
)

type notFound struct {
	engine *Engine
}

func (n *notFound) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	n.engine.notFoundHandler(w, r)
}
