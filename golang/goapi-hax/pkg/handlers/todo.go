package handlers

import (
	"goapi-hax/pkg/core/dto"
	"goapi-hax/pkg/core/ports"
	"net/http"
	"strconv"
)

type todoHandler struct {
	serv ports.TodoService
}

func NewTodoHandler(s ports.TodoService) ports.TodoHandler {
	return &todoHandler{
		serv: s,
	}
}

func (h todoHandler) CreateTodo(c ports.Context) {
	var newTodo dto.NewTodoRequset

	err := c.Bind(&newTodo)
	if err != nil {
		c.Error(err)
		return
	}

	todo, err := h.serv.Create(newTodo)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, todo)
}

func (h todoHandler) ListTodo(c ports.Context) {
	var completed string
	if val, ok := c.Query("completed"); ok {
		completed = val
	}

	todos, err := h.serv.List(completed)

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, todos)
}

func (h todoHandler) GetTodo(c ports.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	todo, err := h.serv.Get(id)

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, todo)
}

func (h todoHandler) UpdateTodo(c ports.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var updateTodo dto.UpdateTodoStatus

	err := c.Bind(&updateTodo)
	if err != nil {
		c.Error(err)
		return
	}

	err = h.serv.Update(id, updateTodo)

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h todoHandler) DeleteTodo(c ports.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	err := h.serv.Delete(id)

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
