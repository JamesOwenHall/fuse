package main

import (
	"github.com/JamesOwenHall/fuse"
)

func main() {
	f := fuse.New()
	f.GET("/", func(c *fuse.Context) {
		c.ResponseWriter.Write([]byte("hello world"))
	})
	f.Run(":3000")
}
