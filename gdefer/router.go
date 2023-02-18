package gdefer

import (
	"log"
	"net/http"
	"strings"
)

//	key := c.Method + "-" + c.Path -> handlers
type Router struct {
	roots    map[string]*node
	handlers map[string]HandleFunc
}

func newRouter() *Router {
	return &Router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandleFunc),
	}
}

func (r *Router) handle(c *Context) {
	n, params := r.getRouter(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		//讓handlers符合handleFunc->handleFunc又為Context對象，將r.handlers轉為context對象
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 not fund path[%s]\n", c.Path)
		})
	}
	c.Next()
}

func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, v := range vs {
		if v != "" {
			parts = append(parts, v)
			if v == "*" {
				break
			}
		}
	}
	return parts
}

func (r *Router) addRouter(method, pattern string, handle HandleFunc) {
	log.Printf("Router : method [%s] - pattern [%s]", method, pattern)
	parts := parsePattern(pattern)
	key := method + "-" + pattern

	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handle
}

func (r *Router) getRouter(method, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string, 0) //存放匹配結果

	if _, ok := r.roots[method]; !ok {
		return nil, nil
	}

	n := r.roots[method].search(searchParts, 0)
	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}
