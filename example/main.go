package main

import (
	"strings"

	"github.com/JamesOwenHall/fuse"
)

func main() {
	f := fuse.New()
	f.Use(fuse.Logger)
	f.NotFound = func(c *fuse.Context) {
		c.ResponseWriter.Write([]byte("Darn, not found"))
	}

	f.GET("/", func(c *fuse.Context) {
		c.TextOk("hello world")
	})
	f.GET("/say", func(c *fuse.Context) {
		words := strings.Join(c.Form["message"], " ")
		c.TextOk(words)
	})
	f.GET("/say/:message", func(c *fuse.Context) {
		c.TextOk(c.Params["message"])
	})

	f.Run(":3000")
}
