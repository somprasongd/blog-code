package server

import (
	"encoding/json"
	"goapi-hax/pkg/common/errs"
	"goapi-hax/pkg/core/ports"
	"net/http"

	"github.com/gorilla/mux"
)

type HttpContext struct {
	w http.ResponseWriter
	r *http.Request
}

func NewHttpContext(w http.ResponseWriter, r *http.Request) ports.Context {
	return &HttpContext{
		w: w,
		r: r,
	}
}

func (c *HttpContext) Bind(v interface{}) error {
	return json.NewDecoder(c.r.Body).Decode(v)
}

func (c *HttpContext) Query(key string) (string, bool) {
	query := c.r.URL.Query()
	if vals, ok := query[key]; ok {
		return vals[0], ok
	}
	return "", false
}

func (c *HttpContext) Param(key string) string {
	vars := mux.Vars(c.r)
	return vars[key]
}

func (c *HttpContext) Error(err error) {
	switch e := err.(type) {
	case errs.AppError:
		c.JSON(e.Code, e)
	case error:
		appErr := errs.AppError{
			Code:    http.StatusInternalServerError,
			Message: e.Error(),
		}
		c.JSON(appErr.Code, appErr)
	}
}

func (c *HttpContext) JSON(code int, data interface{}) {
	c.w.Header().Add("Content-Type", "application/json")
	c.w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(c.w).Encode(data)
	}
}

func NewHttpHandler(handler func(ports.Context)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(NewHttpContext(w, r))
	}
}
