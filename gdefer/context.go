package gdefer

import (
	"encoding/json"
	"fmt"
	"math"
	"net"
	"net/http"
	"strings"
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

	engine *Engine
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

func (c *Context) Html(code int, name string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	if err := c.engine.htmlTemplate.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.Status(http.StatusInternalServerError)
		fmt.Fprintf(c.Writer, "error load html")
		return
	}
}

func (c *Context) GetIP() (ip string, port string, forward string) {
	ip, port, _ = net.SplitHostPort(c.Request.RemoteAddr)
	if ip == "" {
		fmt.Fprintf(c.Writer, "Error for getting IP address")
		return
	}
	userIP := net.ParseIP(ip)
	if userIP == nil {
		fmt.Fprintf(c.Writer, "Error for getting user IP address")
		return
	}
	forWord := c.Request.Header.Get("X-Forwarded-For")
	fmt.Fprintf(c.Writer, "<p>IP: %s</p>", ip)
	fmt.Fprintf(c.Writer, "<p>Port: %s</p>", port)
	fmt.Fprintf(c.Writer, "<p>Forwarded for: %s</p>", forWord)

	xForward := c.Request.Header.Get("X-Forwarded-For")
	Cip := strings.TrimSpace(strings.Split(xForward, ",")[0])
	if ip != "" {
		return Cip, "", ""
	}
	Cip = strings.TrimSpace(strings.Split(c.Request.Header.Get("X-Real-Ip"), ",")[0])
	if ip != "" {
		return Cip, "", ""
	}
	return "", "", ""
}

func (c *Context) Fail(statusCode int, msg string) {
	c.Json(statusCode, msg)
}
