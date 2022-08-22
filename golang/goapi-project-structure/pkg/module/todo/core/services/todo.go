package services

import (
	"errors"
	"goapi/pkg/common"
	"goapi/pkg/common/logger"
	"goapi/pkg/module/todo/core/dto"
	"goapi/pkg/module/todo/core/mapper"
	"goapi/pkg/module/todo/core/model"
	"goapi/pkg/module/todo/core/ports"
)

var (
	// ErrTodoNotFoundById todo not found error when find with id
	ErrTodoNotFoundById = common.NewNotFoundError("todo with given id not found")
)

type todoService struct {
	repo ports.TodoRepository
}

func NewTodoService(repo ports.TodoRepository) ports.TodoService {
	return &todoService{repo}
}

func (s todoService) Create(form dto.NewTodoForm, reqId string) (*dto.TodoResponse, error) {
	// validate
	if err := common.ValidateDto(form); err != nil {
		return nil, common.NewInvalidError(err.Error())
	}

	todo := model.Todo{
		Text: form.Text,
	}
	// เรียกใช้ repo เพื่อบันทึกข้อมูลใหม่
	err := s.repo.Create(&todo)
	if err != nil {
		return nil, common.ErrDbInsert
	}

	// สร้าง struct ที่ต้องการให้ handler ส่งกลับไปหา client
	serialized := mapper.TodoToDto(&todo)

	return serialized, nil
}

func (s todoService) List(page common.PagingRequest, filters dto.ListTodoFilter, reqId string) (dto.TodoResponses, *common.PagingResult, error) {
	// validate
	if err := common.ValidateDto(filters); err != nil {
		return nil, nil, common.NewInvalidError(err.Error())
	}

	todos, paging, err := s.repo.Find(page, filters)
	if err != nil {
		logger.ErrorWithReqId(err.Error(), reqId)
		return nil, nil, common.ErrDbQuery
	}

	serialized := mapper.TodosToDto(todos)
	return serialized, paging, nil
}

func (s todoService) Get(id string, reqId string) (*dto.TodoResponse, error) {
	// validate id format
	if err := s.validateId(id); err != nil {
		return nil, common.ErrIdFormat
	}

	task, err := s.repo.FindById(id)
	if err != nil {
		if errors.Is(err, common.ErrRecordNotFound) {
			return nil, ErrTodoNotFoundById
		}
		logger.ErrorWithReqId(err.Error(), reqId)
		return nil, common.ErrDbQuery
	}

	serialized := mapper.TodoToDto(task)
	return serialized, nil
}

func (s todoService) UpdateStatus(id string, form dto.UpdateTodoForm, reqId string) (*dto.TodoResponse, error) {
	// validate id format
	if err := s.validateId(id); err != nil {
		return nil, common.ErrIdFormat
	}

	err := common.ValidateDto(form)
	if err != nil {
		return nil, common.NewInvalidError(err.Error())
	}

	todo, err := s.repo.UpdateStatusById(id, form.Completed)
	if err != nil {
		if errors.Is(err, common.ErrRecordNotFound) {
			return nil, ErrTodoNotFoundById
		}
		logger.ErrorWithReqId(err.Error(), reqId)
		return nil, common.ErrDbUpdate
	}

	serialized := mapper.TodoToDto(todo)
	return serialized, nil
}

func (s todoService) Delete(id string, reqId string) error {
	// validate id format
	if err := s.validateId(id); err != nil {
		return common.ErrIdFormat
	}

	err := s.repo.DeleteById(id)
	if err != nil {
		if errors.Is(err, common.ErrRecordNotFound) {
			return ErrTodoNotFoundById
		}
		logger.ErrorWithReqId(err.Error(), reqId)
		return common.ErrDbDelete
	}

	return nil
}

func (s todoService) validateId(id string) error {
	return common.ValidateDto(&dto.TodoId{ID: id})
}
