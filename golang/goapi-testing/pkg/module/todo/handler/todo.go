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

func (h TodoHandler) GetTodo(c common.HContext) error {
	id := c.Param("id")

	todo, err := h.serv.Get(id, c.RequestId())

	if err != nil {
		return common.ResponseError(c, err)
	}

	return common.ResponseOk(c, "todo", todo)
}

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

func (h TodoHandler) DeleteTodo(c common.HContext) error {
	id := c.Param("id")

	err := h.serv.Delete(id, c.RequestId())

	if err != nil {
		return common.ResponseError(c, err)
	}

	return common.ResponseNoContent(c)
}
