package gdefer

import (
	"log"
	"net/http"
)

type Router struct {
	routerHandle map[string]HandleFunc
}

func newRouter() *Router {
	return &Router{
		routerHandle: make(map[string]HandleFunc)}
}

func (r *Router) addRouter(method, path string, handle HandleFunc) {
	log.Printf("Router : method [%s] - path [%s]", method, path)
	key := method + "-" + path
	r.routerHandle[key] = handle
}

func (r *Router) handler(c *Context) {
	key := c.Method + "-" + c.Path
	if handler, ok := r.routerHandle[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "page not found 404 ", c.Path)
	}
}
