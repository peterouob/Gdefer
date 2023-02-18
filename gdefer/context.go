package gdefer

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	//original object
	Writer  http.ResponseWriter
	Request *http.Request
	//request object
	Path   string
	Method string
	Params map[string]string
	//response object
	StatusCode int
	//middleware
	handlers []HandleFunc
	index    int8
}

const abortIndex int8 = math.MaxInt8 >> 1

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer:  w,
		Request: r,
		Path:    r.URL.Path,
		Method:  r.Method,
		index:   -1,
	}
}

func (c *Context) Next() {
	c.index++
	length := len(c.handlers)
	for ; c.index < int8(length); c.index++ {
		c.handlers[c.index](c)
	}
}

func (c *Context) Abort() {
	c.index = abortIndex
	//fmt.Println(c.index)
}

func (c *Context) IsAbort() bool {
	return c.index >= abortIndex
}

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

func (c *Context) PostForm(key string) string {
	return c.Request.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	if _, err := c.Writer.Write([]byte(fmt.Sprintf(format, values...))); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}

func (c *Context) Json(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	if _, err := c.Writer.Write(data); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *Context) Html(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	if _, err := c.Writer.Write([]byte(html)); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}

func (c *Context) ClientIP(r *http.Request) string {
	xForward := r.Header.Get("X-Forwarded-For")

	//ip := strings.TrimSpace(strings.Split(xForward, ",")[0])
	//if ip != "" {
	//	return ip
	//}
	//ip = strings.TrimSpace(strings.Split(c.Request.Header.Get("X-Real-Ip"), ",")[0])
	//if ip != "" {
	//	return ip
	//}
	//if ip, _, err := net.SplitHostPort(strings.TrimSpace(c.Request.RemoteAddr)); err != nil {
	//	return ip
	//}
	return xForward
}
