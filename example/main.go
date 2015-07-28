package main

import (
	"strings"

	"github.com/JamesOwenHall/fuse"
)

func main() {
	f := fuse.New([]byte(`my-secret`))
	f.Use(fuse.Logger)
	f.NotFound = func(c *fuse.Context) {
		c.ResponseWriter.Write([]byte("Darn, not found"))
	}

	f.GET("/", func(c *fuse.Context) {
		visits, exists := c.Session["visits"]
		if !exists {
			c.OutData["visits"] = 0
			c.Session["visits"] = 0
		} else {
			numVisits := visits.(int) + 1
			c.OutData["visits"] = numVisits
			c.Session["visits"] = numVisits
		}

		c.HtmlOk("home.tpl")
	})
	f.GET("/say", func(c *fuse.Context) {
		words := strings.Join(c.Form["message"], " ")
		c.TextOk(words)
	})
	f.GET("/say/:message", func(c *fuse.Context) {
		c.OutData["message"] = c.Params["message"]
		c.JsonOk()
	})

	f.Run(":3000")
}
