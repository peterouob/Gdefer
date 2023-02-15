package main

import (
	"gdefer/gdefer"
	"net/http"
)

func main() {
	g := gdefer.New()
	g.Get("/", func(c *gdefer.Context) {
		c.Html(http.StatusOK, "<h1>Hello Gee</h1>")
	})
	g.Get("/hello/:name", func(c *gdefer.Context) {
		// expect /hello?name=geektutu
		c.String(http.StatusOK, "hello %s, you're at %s \n", c.Query("name"), c.Path)
	})

	g.Get("/assets/*filepath", func(c *gdefer.Context) {
		c.Json(http.StatusOK, gdefer.H{"filepath": c.Param("filepath")})
	})

	g.Post("/login", func(c *gdefer.Context) {
		c.Json(http.StatusOK, gdefer.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})
	g.Run(":80")
}
