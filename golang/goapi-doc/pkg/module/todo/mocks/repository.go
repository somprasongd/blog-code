package mocks

import (
	"goapi/pkg/common"
	"goapi/pkg/module/todo/core/dto"
	"goapi/pkg/module/todo/core/model"
	"goapi/pkg/module/todo/core/ports"

	"github.com/stretchr/testify/mock"
)

type todoRepositoryMock struct {
	mock.Mock
}

var _ ports.TodoRepository = &todoRepositoryMock{}

func NewTodoRepositoryMock() *todoRepositoryMock {
	return &todoRepositoryMock{}
}

func (m *todoRepositoryMock) Create(t *model.Todo) error {
	args := m.Called(t)
	return args.Error(0)
}

func (m *todoRepositoryMock) Find(page common.PagingRequest, filters dto.ListTodoFilter) (model.Todos, *common.PagingResult, error) {
	args := m.Called(page, filters)

	var r0 model.Todos
	if args.Get(0) != nil {
		r0 = args.Get(0).(model.Todos)
	}

	var r1 *common.PagingResult
	if args.Get(1) != nil {
		r1 = args.Get(1).(*common.PagingResult)
	}

	return r0, r1, args.Error(2)
}

func (m *todoRepositoryMock) FindById(id string) (*model.Todo, error) {
	args := m.Called(id)

	var r0 *model.Todo
	if args.Get(0) != nil {
		r0 = args.Get(0).(*model.Todo)
	}

	return r0, args.Error(1)
}

func (m *todoRepositoryMock) UpdateStatusById(id string, status bool) (*model.Todo, error) {
	args := m.Called(id, status)

	var r0 *model.Todo
	if args.Get(0) != nil {
		r0 = args.Get(0).(*model.Todo)
	}

	return r0, args.Error(1)
}

func (m *todoRepositoryMock) DeleteById(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
