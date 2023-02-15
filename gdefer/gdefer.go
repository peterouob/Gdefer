package gdefer

import (
	"fmt"
	"log"
	"net/http"
)

type HandleFunc func(w http.ResponseWriter, r *http.Request)

type Engine struct {
	router map[string]HandleFunc
}

func New() *Engine {
	return &Engine{router: make(map[string]HandleFunc)}
}

func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

func (e *Engine) addRouter(method, path string, handle HandleFunc) {
	key := method + "-" + path
	e.router[key] = handle
}

func (e *Engine) Get(pattern string, handle HandleFunc) {
	e.addRouter("GET", pattern, handle)
}
func (e *Engine) Post(pattern string, handle HandleFunc) {
	e.addRouter("POST", pattern, handle)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.Method + "-" + r.URL.Path
	if handle, ok := e.router[key]; !ok {
		_, err := fmt.Fprintf(w, "cannot fund request 404")
		if err != nil {
			log.Printf("write error for 404 on request has failed \n err[%s]", err)
		}
	} else {
		handle(w, r)
	}
}
