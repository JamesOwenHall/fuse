package fuse

import (
	"log"
	"time"
)

// Logger is a middleware component that logs each request to the log package's
// output.
func Logger(c *Context) {
	start := time.Now()
	defer func() {
		if r := recover(); r != nil {
			log.Println("PANIC", c.Request.Method, "->", c.Request.URL.String())
			panic(r)
		}
	}()

	c.Next()
	dur := time.Now().Sub(start)
	log.Println(c.Request.Method, "-", dur.String(), "->", c.Request.URL.String())
}
