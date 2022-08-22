package domain

import "goapi-hax/pkg/core/dto"

type Todo struct {
	ID        int
	Title     string
	Completed bool `gorm:"column:is_done"`
}

func (t Todo) ToDto() dto.TodoResponse {
	return dto.TodoResponse{
		ID:          t.ID,
		Text:        t.Title,
		IsCompleted: t.Completed,
	}
}
