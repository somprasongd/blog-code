package context

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type HttpContext struct {
	w http.ResponseWriter
	r *http.Request
}

func NewHttpContext(w http.ResponseWriter, r *http.Request) MyContext {
	return &HttpContext{
		w: w,
		r: r,
	}
}

func (c *HttpContext) Bind(v interface{}) error {
	return json.NewDecoder(c.r.Body).Decode(v)
}

func (c *HttpContext) BindQuery(v interface{}) error {
	querys := c.r.URL.Query()
	m := map[string]string{}
	for k, v := range querys {
		m[k] = v[0]
	}
	jsonStr, err := json.Marshal(m)
	if err != nil {
		return err
	}

	return json.NewDecoder(strings.NewReader(string(jsonStr))).Decode(v)
}

func (c *HttpContext) Query(key string) (string, bool) {
	query := c.r.URL.Query()
	if vals, ok := query[key]; ok {
		return vals[0], true
	}
	return "", false
}

func (c *HttpContext) DefaultQuery(key string, d string) string {
	query := c.r.URL.Query()
	if vals, ok := query[key]; ok {
		return vals[0]
	}
	return d
}

func (c *HttpContext) Param(key string) string {
	vars := mux.Vars(c.r)
	return vars[key]
}

func (c *HttpContext) Header(key string) string {
	return c.r.Header.Get(key)
}

func (c *HttpContext) RequestId() string {
	return c.Header("X-Request-Id")
}

func (c *HttpContext) ResponseError(code int, err string) {
	c.ResponseJSON(code, map[string]string{
		"error": err,
	})
}

func (c *HttpContext) ResponseJSON(code int, data interface{}) {
	c.w.Header().Add("Content-Type", "application/json")
	c.w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(c.w).Encode(data)
	}
}

func WrapHTTPHandler(h func(MyContext)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		h(NewHttpContext(w, r))
	}
}
