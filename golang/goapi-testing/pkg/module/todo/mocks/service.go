package mocks

import (
	"goapi/pkg/common"
	"goapi/pkg/module/todo/core/dto"
	"goapi/pkg/module/todo/core/ports"

	"github.com/stretchr/testify/mock"
)

type taskServiceMock struct {
	mock.Mock
}

var _ ports.TodoService = &taskServiceMock{}

func NewTaskServiceMock() *taskServiceMock {
	return &taskServiceMock{}
}

func (m *taskServiceMock) Create(form dto.NewTodoForm, reqId string) (*dto.TodoResponse, error) {
	args := m.Called(form, reqId)

	var r0 *dto.TodoResponse
	if args.Get(0) != nil {
		r0 = args.Get(0).(*dto.TodoResponse)
	}

	return r0, args.Error(1)
}

func (m *taskServiceMock) List(page common.PagingRequest, filters dto.ListTodoFilter, reqId string) (dto.TodoResponses, *common.PagingResult, error) {
	args := m.Called(page, filters, reqId)

	var r0 dto.TodoResponses
	if args.Get(0) != nil {
		r0 = args.Get(0).(dto.TodoResponses)
	}

	var r1 *common.PagingResult
	if args.Get(1) != nil {
		r1 = args.Get(1).(*common.PagingResult)
	}

	return r0, r1, args.Error(2)
}

func (m *taskServiceMock) Get(id string, reqId string) (*dto.TodoResponse, error) {
	args := m.Called(id, reqId)

	var r0 *dto.TodoResponse
	if args.Get(0) != nil {
		r0 = args.Get(0).(*dto.TodoResponse)
	}

	return r0, args.Error(1)
}

func (m *taskServiceMock) UpdateStatus(id string, updateTodo dto.UpdateTodoForm, reqId string) (*dto.TodoResponse, error) {
	args := m.Called(id, updateTodo, reqId)

	var r0 *dto.TodoResponse
	if args.Get(0) != nil {
		r0 = args.Get(0).(*dto.TodoResponse)
	}

	return r0, args.Error(1)
}

func (m *taskServiceMock) Delete(id string, reqId string) error {
	args := m.Called(id, reqId)
	return args.Error(0)
}
