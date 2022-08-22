package services

import (
	"goapi-hex/pkg/common/errs"
	"goapi-hex/pkg/core/dto"
	"goapi-hex/pkg/core/model"
	"goapi-hex/pkg/core/ports"
)

type todoService struct {
	repo ports.TodoRepository
}

func NewTodoService(repo ports.TodoRepository) ports.TodoService {
	return &todoService{repo}
}

func (s todoService) Create(form dto.NewTodoForm) (*dto.TodoResponse, error) {
	// การตรวจสอบ
	if form.Text == "" {
		return nil, errs.NewBadRequestError("text is required")
	}

	todo := model.Todo{
		Text: form.Text,
	}
	// เรียกใช้ repo เพื่อบันทึกข้อมูลใหม่
	err := s.repo.Create(&todo)
	if err != nil {
		return nil, errs.NewUnexpectedError("database error while insert new todo")
	}

	// สร้าง struct ที่ต้องการให้ handler ส่งกลับไปหา client
	serializedTodo := dto.TodoResponse{
		ID:   todo.ID,
		Text: todo.Text,
		Done: todo.Done,
	}

	return &serializedTodo, nil
}
