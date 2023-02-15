package gdefer

import (
	"net/http"
)

type HandleFunc func(ctx *Context)

type Engine struct {
	router *Router
}

func New() *Engine {
	return &Engine{router: newRouter()}
}

func (e *Engine) addRouter(method, path string, handle HandleFunc) {
	e.router.addRouter(method, path, handle)
}

func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

func (e *Engine) Get(pattern string, handle HandleFunc) {
	e.addRouter("GET", pattern, handle)
}
func (e *Engine) Post(pattern string, handle HandleFunc) {
	e.addRouter("POST", pattern, handle)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := newContext(w, r)
	e.router.handle(c)
}
