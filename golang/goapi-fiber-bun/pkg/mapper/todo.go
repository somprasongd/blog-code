package mapper

import (
	"goapi/pkg/dtos"
	"goapi/pkg/models"
)

func CreateTodoFormToModel(dto *dtos.CreateTodoForm) *models.Todo {
	return &models.Todo{
		Title: dto.Text,
	}
}

func TodoToDto(m *models.Todo) *dtos.TodoResponse {
	return &dtos.TodoResponse{
		ID:        m.ID,
		Text:      m.Title,
		Completed: m.IsDone,
	}
}

func TodosToDto(todos []models.Todo) []*dtos.TodoResponse {
	dtos := make([]*dtos.TodoResponse, len(todos))
	for i, t := range todos {
		dtos[i] = TodoToDto(&t)
	}

	return dtos
}
