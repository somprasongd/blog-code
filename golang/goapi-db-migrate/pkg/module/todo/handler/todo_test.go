package handler_test

import (
	"bytes"
	"encoding/json"
	"goapi/pkg/common"
	"goapi/pkg/module/todo"
	"goapi/pkg/module/todo/core/dto"
	"goapi/pkg/module/todo/mocks"
	"io"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTodoHandler(t *testing.T) {

	t.Run("Add Todo", func(t *testing.T) {

		t.Run("success_200", func(t *testing.T) {
			//Arrange
			mockDto := dto.TodoResponse{
				ID:        "bfbc2a69-9825-4a0e-a8d6-ffb985dc719c",
				Text:      "New todo",
				Completed: false,
			}

			expectedCode := http.StatusCreated
			mockResp := common.Response{
				Status: expectedCode,
				Data: map[string]interface{}{
					"todo": mockDto,
				},
				RequestId: mock.Anything,
			}

			svc := mocks.NewTaskServiceMock()
			svc.On("Create", mock.AnythingOfType("dto.NewTodoForm"), mock.Anything).Return(&mockDto, nil)

			//http://localhost:8000/todos
			router := fiber.New()
			router.Use(requestid.New())
			cfg := todo.RouteConfig{
				Router:      router,
				TodoService: svc,
			}
			todo.SetupRoutes(cfg)

			reqBody, err := json.Marshal(map[string]string{
				"text": "New todo",
			})

			assert.NoError(t, err)

			req, err := http.NewRequest("POST", "/todos", bytes.NewReader(reqBody))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Request-ID", mock.Anything)

			//Act
			resp, err := router.Test(req)
			assert.NoError(t, err)

			//Assert
			if assert.Equal(t, expectedCode, resp.StatusCode) {
				body, _ := io.ReadAll(resp.Body)
				expected, _ := json.Marshal(mockResp)
				assert.JSONEq(t, string(expected), string(body))
			}
		})

		t.Run("bodyparser_error_400", func(t *testing.T) {
			//Arrange
			expectedCode := http.StatusBadRequest
			err := common.ErrBodyParser
			mockResp := common.Response{
				Status:    expectedCode,
				Error:     err,
				RequestId: mock.Anything,
			}

			svc := mocks.NewTaskServiceMock()

			//http://localhost:8000/todos
			router := fiber.New()
			router.Use(requestid.New())
			cfg := todo.RouteConfig{
				Router:      router,
				TodoService: svc,
			}
			todo.SetupRoutes(cfg)

			reqBody, err := json.Marshal(map[string]int{
				"text": 1,
			})

			assert.NoError(t, err)

			req, err := http.NewRequest("POST", "/todos", bytes.NewReader(reqBody))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Request-ID", mock.Anything)
			//Act
			resp, err := router.Test(req)
			assert.NoError(t, err)

			//Assert
			if assert.Equal(t, expectedCode, resp.StatusCode) {
				body, _ := io.ReadAll(resp.Body)
				expected, _ := json.Marshal(mockResp)
				assert.JSONEq(t, string(expected), string(body))
			}
		})

		t.Run("invalid_body_422", func(t *testing.T) {
			//Arrange
			expectedCode := http.StatusUnprocessableEntity
			err := common.NewInvalidError("text: text is a required field")
			mockResp := common.Response{
				Status:    expectedCode,
				Error:     err,
				RequestId: mock.Anything,
			}

			svc := mocks.NewTaskServiceMock()
			svc.On("Create", mock.AnythingOfType("dto.NewTodoForm"), mock.Anything).Return(nil, err)

			//http://localhost:8000/todos
			router := fiber.New()
			router.Use(requestid.New())
			cfg := todo.RouteConfig{
				Router:      router,
				TodoService: svc,
			}
			todo.SetupRoutes(cfg)

			reqBody, err := json.Marshal(map[string]string{
				"text": "",
			})

			assert.NoError(t, err)

			req, err := http.NewRequest("POST", "/todos", bytes.NewReader(reqBody))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Request-ID", mock.Anything)
			//Act
			resp, err := router.Test(req)
			assert.NoError(t, err)

			//Assert
			if assert.Equal(t, expectedCode, resp.StatusCode) {
				body, _ := io.ReadAll(resp.Body)
				expected, _ := json.Marshal(mockResp)
				assert.JSONEq(t, string(expected), string(body))
			}
		})

		t.Run("error_500", func(t *testing.T) {
			//Arrange
			expectedCode := http.StatusInternalServerError
			err := common.ErrDbInsert
			mockResp := common.Response{
				Status:    expectedCode,
				Error:     err,
				RequestId: mock.Anything,
			}

			svc := mocks.NewTaskServiceMock()
			svc.On("Create", mock.AnythingOfType("dto.NewTodoForm"), mock.Anything).Return(nil, err)

			//http://localhost:8000/todos
			router := fiber.New()
			router.Use(requestid.New())
			cfg := todo.RouteConfig{
				Router:      router,
				TodoService: svc,
			}
			todo.SetupRoutes(cfg)

			reqBody, err := json.Marshal(map[string]string{
				"text": "New todo",
			})

			assert.NoError(t, err)

			req, err := http.NewRequest("POST", "/todos", bytes.NewReader(reqBody))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Request-ID", mock.Anything)
			//Act
			resp, err := router.Test(req)
			assert.NoError(t, err)

			//Assert
			if assert.Equal(t, expectedCode, resp.StatusCode) {
				body, _ := io.ReadAll(resp.Body)
				expected, _ := json.Marshal(mockResp)
				assert.JSONEq(t, string(expected), string(body))
			}
		})

	})

}
