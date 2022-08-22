package main

import (
	"goapi-handlefunc/context"
	"goapi-handlefunc/handler"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	BASE_URL = "/api/v1"
	PORT     = ":8080"
)

func main() {
	r := gin.Default()

	setRouter(r)

	http.ListenAndServe(PORT, r)
}

func setRouter(r *gin.Engine) {
	h := handler.TodoHandler{}

	todos := r.Group(BASE_URL + "/todos")
	todos.POST("", context.WrapGinHandler(h.CreateHandler))
	todos.GET("", context.WrapGinHandler(h.ListHandler))
	todos.GET("/:id", context.WrapGinHandler(h.GetHandler))
	todos.PATCH("/:id", context.WrapGinHandler(h.StatusUpdateHandler))
	todos.DELETE("/:id", context.WrapGinHandler(h.DeleteHandler))
}
