package main

import (
	"fmt"
	"gdefer/gdefer"
	"net/http"
)

func onlyForV2() gdefer.HandleFunc {
	return func(c *gdefer.Context) {
		fmt.Println(123)
		c.Abort()
	}
}

func main() {
	g := gdefer.New()
	g.Use(gdefer.Logger())
	g.Get("/hello", func(c *gdefer.Context) {
		// expect /hello?name=geektutu

		c.Json(http.StatusOK, gdefer.H{
			"ip": c.ClientIP(c.Request),
		})
	})
	v1 := g.Group("/v1")
	v1.Use(onlyForV2())
	{
		v1.Get("/", func(c *gdefer.Context) {
			c.Abort()
			c.Json(200, gdefer.H{
				"status": "success",
			})
		})
	}
	g.Run(":80")
}
