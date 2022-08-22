package dto

type NewTodoForm struct {
	Text string `json:"text" validate:"required"`
}

type ListTodoFilter struct {
	Term      string
	Completed *bool
}

type TodoId struct {
	ID string `validate:"required,uuid4"`
}

type UpdateTodoForm struct {
	Completed bool `json:"completed" validate:"required"`
}

type TodoResponse struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	Completed bool   `json:"isCompleted"`
}

type TodoResponses []*TodoResponse
