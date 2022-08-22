package handler

import (
	"fmt"
	"goapi-handlefunc/context"
	"net/http"
	"strconv"
)

type Todo struct {
	ID        int    `json:"id"`
	Text      string `json:"text"`
	Completed bool   `json:"completed"`
}

type TodoHandler struct {
}

func (h TodoHandler) CreateHandler(ctx context.MyContext) {
	// bind json to new stuct
	var todo Todo
	err := ctx.Bind(&todo)
	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return
	}
	// validate stuct if error return error reponse with status 400
	if todo.Text == "" {
		ctx.ResponseError(http.StatusBadRequest, "text is required")
		return
	}
	// save
	// return json response with status 201
	ctx.ResponseJSON(http.StatusCreated, todo)
}

func (h TodoHandler) ListHandler(ctx context.MyContext) {
	// get query param for filter
	term := "Text"
	if val, ok := ctx.Query("term"); ok {
		term = val
	}
	// list by filter
	todos := []Todo{
		{ID: 1, Text: term + "1", Completed: true},
		{ID: 2, Text: term + "2", Completed: false},
		{ID: 3, Text: term + "3", Completed: false},
	}
	// return json response with status 200
	ctx.ResponseJSON(http.StatusOK, todos)
}

func (h TodoHandler) GetHandler(ctx context.MyContext) {
	// get id from path param
	id, _ := strconv.Atoi(ctx.Param("id"))
	// get by id
	todo := Todo{
		ID:        id,
		Text:      "Get Todo by ID",
		Completed: false,
	}
	// return json notfound error reponse if notfound with status 404
	// return json response with status 200
	ctx.ResponseJSON(http.StatusOK, todo)
}

func (h TodoHandler) StatusUpdateHandler(ctx context.MyContext) {
	// get id from path param
	id, _ := strconv.Atoi(ctx.Param("id"))
	// bind json to patch stuct
	var updateTodo Todo
	err := ctx.Bind(&updateTodo)
	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return
	}
	// get by id
	todo := Todo{
		ID:        id,
		Text:      "Update Todo Status by ID",
		Completed: false,
	}
	// return json notfound error reponse if notfound with status 404
	// udpate status
	todo.Completed = updateTodo.Completed
	// return json response with status 200
	ctx.ResponseJSON(http.StatusOK, todo)
}

func (h TodoHandler) DeleteHandler(ctx context.MyContext) {
	// get id from path param
	id, _ := strconv.Atoi(ctx.Param("id"))
	// get by id
	// return json notfound error reponse if notfound with status 404
	// delete by id
	fmt.Println("Delete Todo by ID:", id)
	// return empty response with status 204
	ctx.ResponseJSON(http.StatusNoContent, nil)
}
