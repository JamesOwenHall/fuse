package fuse

import (
	"log"
	"time"
)

func Logger(c *Context) {
	start := time.Now()
	c.Next()
	dur := time.Now().Sub(start)
	log.Println(c.Request.Method, "-", dur.String(), "->", c.Request.URL.String())
}
