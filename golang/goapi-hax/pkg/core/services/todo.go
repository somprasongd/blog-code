package services

import (
	"goapi-hax/pkg/common/errs"
	"goapi-hax/pkg/common/logger"
	"goapi-hax/pkg/common/validator"
	"goapi-hax/pkg/core/domain"
	"goapi-hax/pkg/core/dto"
	"goapi-hax/pkg/core/ports"
	"strconv"
)

type todoService struct {
	repo ports.TodoRepository
}

func NewTodoService(r ports.TodoRepository) ports.TodoService {
	return &todoService{
		repo: r,
	}
}

func (s todoService) Create(newTodo dto.NewTodoRequset) (*dto.TodoResponse, error) {
	// validate
	err := validator.ValidateStruct(newTodo)
	if err != nil {
		return nil, errs.NewBadRequestError(err.Error())
	}
	// create new domain
	todo := domain.Todo{
		Title: newTodo.Text,
	}
	// save to db
	err = s.repo.Create(&todo)
	if err != nil {
		logger.Error(err.Error())
		return nil, errs.NewUnexpectedError(err.Error())
	}
	// convert to dto
	todoResp := todo.ToDto()
	return &todoResp, nil
}
func (s todoService) List(completed string) ([]dto.TodoResponse, error) {
	conditions := map[string]interface{}{}
	if completed != "" {
		b1, err := strconv.ParseBool(completed)
		if err != nil {
			return nil, errs.NewBadRequestError("completed must is boolean")
		}
		conditions["is_done"] = b1
	}
	todos, err := s.repo.Find(conditions)
	if err != nil {
		logger.Error(err.Error())
		return nil, errs.NewUnexpectedError(err.Error())
	}
	// convert to dto
	todoResps := []dto.TodoResponse{}
	for _, t := range todos {
		todoResps = append(todoResps, t.ToDto())
	}
	return todoResps, nil
}
func (s todoService) Get(id int) (*dto.TodoResponse, error) {
	todo, err := s.repo.FindById(id)
	if err != nil {
		if err.Error() == "NOT_FOUND" {
			return nil, errs.NewNotFoundError("todo with given id not found")
		}
		logger.Error(err.Error())
		return nil, errs.NewUnexpectedError(err.Error())
	}
	todoResp := todo.ToDto()
	return &todoResp, nil
}
func (s todoService) Update(id int, updateTodo dto.UpdateTodoStatus) error {
	// validate
	err := validator.ValidateStruct(updateTodo)
	if err != nil {
		return errs.NewBadRequestError(err.Error())
	}
	err = s.repo.UpdateStatusById(id, updateTodo.IsCompleted)
	if err != nil {
		if err.Error() == "NOT_FOUND" {
			return errs.NewNotFoundError("todo with given id not found")
		}
		logger.Error(err.Error())
		return errs.NewUnexpectedError(err.Error())
	}
	return nil
}
func (s todoService) Delete(id int) error {
	err := s.repo.DeleteById(id)
	if err != nil {
		if err.Error() == "NOT_FOUND" {
			return errs.NewNotFoundError("todo with given id not found")
		}
		logger.Error(err.Error())
		return errs.NewUnexpectedError(err.Error())
	}
	return nil
}
