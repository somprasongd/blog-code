package handler

import (
	"goapi/pkg/common"
	"goapi/pkg/module/todo/core/dto"
	"goapi/pkg/module/todo/core/ports"
)

type TodoHandler struct {
	serv ports.TodoService
}

func NewTodoHandler(serv ports.TodoService) *TodoHandler {
	return &TodoHandler{serv}
}

// @Summary Add a new todo
// @Description Add a new todo
// @Tags Todo
// @Accept  json
// @Produce  json
// @Param todo body swagger.CreateTodoFrom true "Todo Data"
// @Failure 422 {object} swagdto.Error422{error=swagger.ErrCreateSampleData}
// @Failure 500 {object} swagdto.Error500
// @Success 201 {object} swagdto.Response{data=swagger.TodoSampleData}
// @Router /todos [post]
func (h TodoHandler) CreateTodo(c common.HContext) error {
	// แปลง JSON เป็น struct
	form := new(dto.NewTodoForm)
	if err := c.BodyParser(form); err != nil {
		return common.ResponseError(c, common.ErrBodyParser)
	}
	// ส่งต่อไปให้ service ทำงาน
	todo, err := h.serv.Create(*form, c.RequestId())
	if err != nil {
		// error จะถูกจัดการมาจาก service แล้ว
		return common.ResponseError(c, err)
	}

	// คืนค่า todo ที่เพิ่งบันทึกเสร็จกลับไปในรูปแบบ JSON
	return common.ResponseCreated(c, "todo", todo)
}

// @Summary List all existing todos
// @Description You can filter all existing todos by listing them.
// @Tags Todo
// @Accept  json
// @Produce  json
// @Param term query string false "filter the text based value (ex: term=dosomething)"
// @Param completed query bool false "filter the status based value (ex: completed=true)"
// @Param page query int false "Go to a specific page number. Start with 1"
// @Param limit query int false "Page size for the data"
// @Param order query string false "Page order. Eg: text desc,createdAt desc"
// @Failure 400 {object} swagdto.Error400
// @Failure 500 {object} swagdto.Error500
// @Success 200 {object} swagdto.ResponseWithPage{data=swagger.TodoSampleListData}
// @Router /todos [get]
func (h TodoHandler) ListTodo(c common.HContext) error {
	filters := dto.ListTodoFilter{}
	if err := c.QueryParser(&filters); err != nil {
		return common.ResponseError(c, common.ErrQueryParser)
	}

	page := common.Paginator(c)

	todos, paging, err := h.serv.List(page, filters, c.RequestId())

	if err != nil {
		return common.ResponseError(c, err)
	}

	return common.ResponsePage(c, "todos", todos, paging)
}

// @Summary Get a todo
// @Description Get a specific todo by id
// @Produce json
// @Tags Todo
// @Param id path string true "Todo ID"
// @Failure 400 {object} swagdto.Error400
// @Failure 404 {object} swagdto.Error404
// @Failure 500 {object} swagdto.Error500
// @Success 200 {object} swagdto.Response{data=swagger.TodoSampleData}
// @Router /todos/{id} [get]
func (h TodoHandler) GetTodo(c common.HContext) error {
	id := c.Param("id")

	todo, err := h.serv.Get(id, c.RequestId())

	if err != nil {
		return common.ResponseError(c, err)
	}

	return common.ResponseOk(c, "todo", todo)
}

// @Summary Update a todo status
// @Description Update a specific todo status by id
// @Produce json
// @Tags Todo
// @Param id path string true "Todo ID"
// @Param todo body swagger.UpdateTodoStatusForm true "Todo Status Data"
// @Failure 400 {object} swagdto.Error400
// @Failure 404 {object} swagdto.Error404
// @Failure 422 {object} swagdto.Error422{error=swagger.ErrUpdateSampleData}
// @Failure 500 {object} swagdto.Error500
// @Success 200 {object} swagdto.Response{data=swagger.TodoSampleData}
// @Router /todos/{id} [patch]
func (h TodoHandler) UpdateTodoStatus(c common.HContext) error {
	id := c.Param("id")

	form := dto.UpdateTodoForm{}

	if err := c.BodyParser(&form); err != nil {
		return common.ResponseError(c, err)
	}

	todo, err := h.serv.UpdateStatus(id, form, c.RequestId())

	if err != nil {
		return common.ResponseError(c, err)
	}

	return common.ResponseOk(c, "todo", todo)
}

// @Summary Delete a todo
// @Description Delete a specific todo by id
// @Produce  json
// @Tags Todo
// @Param id path string true "Todo ID"
// @Failure 400 {object} swagdto.Error400
// @Failure 404 {object} swagdto.Error404
// @Failure 500 {object} swagdto.Error500
// @Success 204
// @Router /todos/{id} [delete]
func (h TodoHandler) DeleteTodo(c common.HContext) error {
	id := c.Param("id")

	err := h.serv.Delete(id, c.RequestId())

	if err != nil {
		return common.ResponseError(c, err)
	}

	return common.ResponseNoContent(c)
}
