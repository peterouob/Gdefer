package gdefer

import (
	"net/http"
)

type HandleFunc func(ctx *Context)

type RouterGroup struct {
	prefix     string
	middleware []HandleFunc
	parent     *RouterGroup
	engine     *Engine
}

type Engine struct {
	*RouterGroup
	router *Router
	groups []*RouterGroup
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (g *RouterGroup) Group(prefix string) *RouterGroup {
	engine := g.engine
	newGroup := &RouterGroup{
		prefix: g.prefix + prefix,
		parent: g,
		engine: engine,
	}
	return newGroup
}

func (g *RouterGroup) addRouter(method, comb string, handle HandleFunc) {
	pattern := g.prefix + comb
	//log.Printf("method [%s] pattern [%s]", method, pattern)
	g.engine.router.addRouter(method, pattern, handle)
}

func (g *RouterGroup) Get(pattern string, handle HandleFunc) {
	g.addRouter("GET", pattern, handle)
}
func (g *RouterGroup) Post(pattern string, handle HandleFunc) {
	g.addRouter("POST", pattern, handle)
}

func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := newContext(w, r)
	e.router.handle(c)
}
