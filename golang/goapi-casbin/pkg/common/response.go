package common

import (
	"net/http"
)

type Response struct {
	Status    int                    `json:"status" example:"200"`
	Data      map[string]interface{} `json:"data,omitempty" example:"{data:{task}}"`
	Error     interface{}            `json:"error,omitempty" example:"{}"`
	RequestId string                 `json:"requestId" example:"3b6272b9-1ef1-45e0"`
}

type ResponseWithPage struct {
	Status     int                    `json:"status" example:"200"`
	Data       map[string]interface{} `json:"data,omitempty" example:"{data:{tasks}}"`
	Error      interface{}            `json:"error,omitempty" example:"{}"`
	Pagination interface{}            `json:"_pagination,omitempty" example:"{}"`
	RequestId  string                 `json:"requestId" example:"3b6272b9-1ef1-45e0"`
}

func ResponseOk(c HContext, key string, body interface{}) error {
	if body == nil {
		return c.SendStatus(http.StatusOK)
	}
	return c.SendJSON(http.StatusOK, Response{
		Status:    http.StatusOK,
		Data:      map[string]interface{}{key: body},
		Error:     nil,
		RequestId: c.RequestId(),
	})
}

func ResponseCreated(c HContext, key string, body interface{}) error {
	if body == nil {
		return c.SendStatus(http.StatusCreated)
	}
	res := Response{
		Status:    http.StatusCreated,
		Data:      map[string]interface{}{key: body},
		Error:     nil,
		RequestId: c.RequestId(),
	}

	return c.SendJSON(http.StatusCreated, res)
}

func ResponseNoContent(c HContext) error {
	return c.SendStatus(http.StatusNoContent)
}

func ResponsePage(c HContext, key string, body interface{}, page interface{}) error {
	dataWithPage := ResponseWithPage{
		Status:     http.StatusOK,
		Data:       map[string]interface{}{key: body},
		Error:      nil,
		Pagination: page,
		RequestId:  c.RequestId(),
	}

	return c.SendJSON(http.StatusOK, dataWithPage)
}

func ResponseError(c HContext, err error) error {
	var appErr AppError
	switch e := err.(type) {
	case AppError:
		appErr = e
	default: // case error
		appErr = AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	res := Response{
		Status:    appErr.Code,
		Data:      nil,
		Error:     appErr,
		RequestId: c.RequestId(),
	}
	return c.SendJSON(appErr.Code, res)
}
