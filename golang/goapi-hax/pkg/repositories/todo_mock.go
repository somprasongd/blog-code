package repositories

import (
	"goapi-hax/pkg/core/domain"
	"goapi-hax/pkg/core/ports"
)

type todoRepositoryMock struct {
	todos []domain.Todo
}

func NewTodoRepositoryMock() ports.TodoRepository {
	return &todoRepositoryMock{
		todos: []domain.Todo{},
	}
}

func (r *todoRepositoryMock) Create(t *domain.Todo) error {
	t.ID = len(r.todos) + 1
	r.todos = append(r.todos, *t)
	return nil
}

func (r *todoRepositoryMock) Find(conditions map[string]interface{}) ([]domain.Todo, error) {
	return r.todos, nil
}

func (r *todoRepositoryMock) FindById(id int) (*domain.Todo, error) {
	for _, todo := range r.todos {
		if todo.ID == id {
			return &todo, nil
		}
	}
	return nil, nil
}

func (r *todoRepositoryMock) UpdateStatusById(id int, isDone bool) error {
	todos := &r.todos
	for index, todo := range *todos {
		if todo.ID == id {
			(*todos)[index].Completed = isDone
			break
		}
	}
	return nil
}
func (r *todoRepositoryMock) DeleteById(id int) error {
	todos := r.todos
	for i, todo := range todos {
		if todo.ID == id {
			r.todos = append(todos[:i], todos[i+1:]...)
			break
		}
	}
	return nil
}
