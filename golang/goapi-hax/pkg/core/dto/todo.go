package dto

type NewTodoRequset struct {
	Text string `json:"text" validate:"required"`
}

type UpdateTodoStatus struct {
	IsCompleted bool `json:"isCompleted" validate:"required"`
}

type TodoResponse struct {
	ID          int    `json:"id"`
	Text        string `json:"text"`
	IsCompleted bool   `json:"isCompleted"`
}
