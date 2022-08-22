package handlers

import (
	"goapi-hex/pkg/common/errs"
	"goapi-hex/pkg/core/dto"
	"goapi-hex/pkg/core/ports"

	"github.com/gofiber/fiber/v2"
)

type TodoHandler struct {
	serv ports.TodoService
}

func NewTodoHandler(serv ports.TodoService) *TodoHandler {
	return &TodoHandler{serv}
}

func (h TodoHandler) CreateTodo(c *fiber.Ctx) error {
	// แปลง JSON เป็น struct
	todoForm := new(dto.NewTodoForm)
	if err := c.BodyParser(todoForm); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})

	}
	// ส่งต่อไปให้ service ทำงาน
	todo, err := h.serv.Create(*todoForm)
	if err != nil {
		// error จะถูกจัดการมาจาก service แล้ว
		appErr := err.(errs.AppError)
		return c.Status(appErr.Code).JSON(appErr)
	}

	// คืนค่า todo ที่เพิ่งบันทึกเสร็จกลับไปในรูปแบบ JSON
	return c.JSON(todo)
}

func (h TodoHandler) ListTodo(c *fiber.Ctx) error {
	return c.SendString("List Todo")
}

func (h TodoHandler) GetTodo(c *fiber.Ctx) error {
	return c.SendString("Get Todo")
}

func (h TodoHandler) UpdateTodo(c *fiber.Ctx) error {
	return c.SendString("Update Todo")
}

func (h TodoHandler) DeleteTodo(c *fiber.Ctx) error {
	return c.SendString("Delete Todo")
}
