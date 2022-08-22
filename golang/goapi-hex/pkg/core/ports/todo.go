package ports

import (
	"goapi-hex/pkg/core/dto"
	"goapi-hex/pkg/core/model"
)

// interface สำหรับ output port
type TodoRepository interface {
	Create(t *model.Todo) error
	// Find() ([]model.Todo, error)
	// FindById(id int) (*model.Todo, error)
	// UpdateStatusById(id int, isDone bool) error
	// DeleteById(id int) error
}

// interface สำหรับ input port
type TodoService interface {
	Create(newTodo dto.NewTodoForm) (*dto.TodoResponse, error)
	// List(completed string) ([]dto.TodoResponse, error)
	// Get(id int) (*dto.TodoResponse, error)
	// Update(id int, updateTodo dto.UpdateTodoForm) error
	// Delete(id int) error
}
