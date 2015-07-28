package fuse

import (
	"log"
	"time"
)

// Logger is a middleware component that logs each request to the log package's
// output.
func Logger(c *Context) {
	start := time.Now()
	c.Next()
	dur := time.Now().Sub(start)
	log.Println(c.Request.Method, "-", dur.String(), "->", c.Request.URL.String())
}
