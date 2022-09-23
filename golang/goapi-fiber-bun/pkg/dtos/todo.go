package dtos

type CreateTodoForm struct {
	Text string `json:"text" validate:"required"`
}

type UpdateTodoForm struct {
	Completed bool `json:"completed" validate:"required"`
}

type TodoResponse struct {
	ID        int    `json:"id"`
	Text      string `json:"text"`
	Completed bool   `json:"completed"`
}
