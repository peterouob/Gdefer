package gdefer

import (
	"log"
	"time"
)

func Logger() HandleFunc {
	return func(c *Context) {
		t := time.Now()
		c.Next()
		log.Printf("[%d] %s in %v", c.StatusCode, c.Request.URL, t)
	}
}
