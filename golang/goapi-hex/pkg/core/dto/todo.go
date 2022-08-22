package dto

type NewTodoForm struct {
	Text string `json:"text"`
}

type UpdateTodoForm struct {
	Done bool `json:"done"`
}

type TodoResponse struct {
	ID   uint   `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}
