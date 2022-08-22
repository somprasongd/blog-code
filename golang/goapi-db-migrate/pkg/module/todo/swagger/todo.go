package swagger

type TodoRepsonse struct {
	ID        string `json:"id" example:"bfbc2a69-9825-4a0e-a8d6-ffb985dc719c"`
	Text      string `json:"text" example:"do something"`
	Completed bool   `json:"completed" example:"false"`
}

type ListTodoRepsonse []TodoRepsonse

type TodoSampleData struct {
	Data TodoRepsonse `json:"todo"`
}

type TodoSampleListData struct {
	Data ListTodoRepsonse `json:"todos"`
}

type CreateTodoFrom struct {
	// Required: true
	Text string `json:"text" example:"do something"`
}

type UpdateTodoStatusForm struct {
	// Required: true
	Completed bool `json:"completed"`
}

type ErrorDetailCreate struct {
	Target  string `json:"target" example:"text"`
	Message string `json:"message" example:"text field is required"`
}

type ErrCreateSampleData struct {
	Code    string              `json:"code" example:"422"`
	Message string              `json:"message" example:"invalid data see details"`
	Details []ErrorDetailCreate `json:"details"`
}

type ErrorDetailUpdate struct {
	Target  string `json:"target" example:"completed"`
	Message string `json:"message" example:"completed field is required"`
}

type ErrUpdateSampleData struct {
	Code    string              `json:"code" example:"422"`
	Message string              `json:"message" example:"invalid data see details"`
	Details []ErrorDetailUpdate `json:"details"`
}
