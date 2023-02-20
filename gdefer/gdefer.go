package gdefer

import (
	"fmt"
	"html/template"
	"net/http"
	"path"
	"strings"
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
	//static
	htmlTemplate *template.Template
	funcMap      template.FuncMap
}

func Default() *Engine {
	engine := New()
	engine.Use(Logger(), Recovery())
	return engine
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (g *RouterGroup) Use(middleware ...HandleFunc) {
	g.middleware = append(g.middleware, middleware...)
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

func (g *RouterGroup) createStaticHandle(relativePath string, fs http.FileSystem) HandleFunc {
	absolutePath := path.Join(g.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))

	return func(c *Context) {
		file := c.Param("file")
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Request)
	}
}

func (g *RouterGroup) Static(relativePath, root string) {
	handler := g.createStaticHandle(relativePath, http.Dir(root))
	urlPath := path.Join(relativePath, "/*filepath")
	g.Get(urlPath, handler)
}

func (e *Engine) setFuncMap(funcMap template.FuncMap) {
	e.funcMap = funcMap
}

func (e *Engine) LoadHtmlGlob(pattern string) {
	e.htmlTemplate = template.Must(template.New("").Funcs(e.funcMap).Parse(pattern))
}

func (e *Engine) Run(addr string) error {
	flag := strings.HasPrefix(addr, ":")
	if !flag {
		fmt.Println("addr must need : ")
		return nil
	}
	return http.ListenAndServe(addr, e)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var middleware []HandleFunc
	for _, groups := range e.groups {
		if strings.HasPrefix(r.URL.Path, groups.prefix) {
			middleware = append(middleware, groups.middleware...)
		}
	}
	c := newContext(w, r)
	c.handlers = middleware
	c.engine = e
	e.router.handle(c)
}
