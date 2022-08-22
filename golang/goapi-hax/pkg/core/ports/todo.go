package ports

import (
	"goapi-hax/pkg/core/domain"
	"goapi-hax/pkg/core/dto"
)

type TodoRepository interface {
	Create(t *domain.Todo) error
	Find(conditions map[string]interface{}) ([]domain.Todo, error)
	FindById(id int) (*domain.Todo, error)
	UpdateStatusById(id int, isDone bool) error
	DeleteById(id int) error
}

type TodoService interface {
	Create(newTodo dto.NewTodoRequset) (*dto.TodoResponse, error)
	List(completed string) ([]dto.TodoResponse, error)
	Get(id int) (*dto.TodoResponse, error)
	Update(id int, updateTodo dto.UpdateTodoStatus) error
	Delete(id int) error
}

type TodoHandler interface {
	CreateTodo(c Context)
	ListTodo(c Context)
	GetTodo(c Context)
	UpdateTodo(c Context)
	DeleteTodo(c Context)
}
