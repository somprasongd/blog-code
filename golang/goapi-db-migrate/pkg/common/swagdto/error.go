package swagdto

type ErrorDetail struct {
	Target  string `json:"target" example:"name"`
	Message string `json:"message" example:"name field is required"`
}

type ErrorData400 struct {
	Code    string `json:"code" example:"400"`
	Message string `json:"message" example:"Bad Request"`
}

type Error400 struct {
	Status    uint         `json:"status" example:"400"`
	Error     ErrorData400 `json:"error"`
	RequestId string       `json:"requestId" example:"3b6272b9-1ef1-45e0"`
}

type ErrorData401 struct {
	Code    string `json:"code" example:"401"`
	Message string `json:"message" example:"Unauthorized"`
}

type Error401 struct {
	Status    uint         `json:"status" example:"401"`
	Error     ErrorData401 `json:"error"`
	RequestId string       `json:"requestId" example:"3b6272b9-1ef1-45e0"`
}

type ErrorData403 struct {
	Code    string `json:"code" example:"403"`
	Message string `json:"message" example:"Forbidden"`
}

type Error403 struct {
	Status    uint         `json:"status" example:"403"`
	Error     ErrorData403 `json:"error"`
	RequestId string       `json:"requestId" example:"3b6272b9-1ef1-45e0"`
}

type ErrorData404 struct {
	Code    string `json:"code" example:"404"`
	Message string `json:"message" example:"Not Found"`
}

type Error404 struct {
	Status    uint         `json:"status" example:"404"`
	Error     ErrorData404 `json:"error"`
	RequestId string       `json:"requestId" example:"3b6272b9-1ef1-45e0"`
}

type ErrorData422 struct {
	Code    string        `json:"code" example:"422"`
	Message string        `json:"message" example:"invalid data see details"`
	Details []ErrorDetail `json:"details"`
}

type Error422 struct {
	Status    uint         `json:"status" example:"422"`
	Error     ErrorData422 `json:"error"`
	RequestId string       `json:"requestId" example:"3b6272b9-1ef1-45e0"`
}

type ErrorData500 struct {
	Code    string `json:"code" example:"500"`
	Message string `json:"message" example:"Internal Server Error"`
}

type Error500 struct {
	Status    uint         `json:"status" example:"500"`
	Error     ErrorData500 `json:"error"`
	RequestId string       `json:"requestId" example:"3b6272b9-1ef1-45e0"`
}
