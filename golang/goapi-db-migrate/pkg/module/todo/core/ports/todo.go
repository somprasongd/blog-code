package ports

import (
	"goapi/pkg/common"
	"goapi/pkg/module/todo/core/dto"
	"goapi/pkg/module/todo/core/model"
)

// interface สำหรับ output port
type TodoRepository interface {
	Create(m *model.Todo) error
	Find(page common.PagingRequest, filters dto.ListTodoFilter) (model.Todos, *common.PagingResult, error)
	FindById(id string) (*model.Todo, error)
	UpdateStatusById(id string, status bool) (*model.Todo, error)
	DeleteById(id string) error
}

// interface สำหรับ input port
type TodoService interface {
	Create(newTodo dto.NewTodoForm, reqId string) (*dto.TodoResponse, error)
	List(page common.PagingRequest, filters dto.ListTodoFilter, reqId string) (dto.TodoResponses, *common.PagingResult, error)
	Get(id string, reqId string) (*dto.TodoResponse, error)
	UpdateStatus(id string, updateTodo dto.UpdateTodoForm, reqId string) (*dto.TodoResponse, error)
	Delete(id string, reqId string) error
}
