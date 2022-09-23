package handlers

import (
	"database/sql"
	"goapi/pkg/common/validator"
	"goapi/pkg/dtos"
	"goapi/pkg/mapper"
	"goapi/pkg/models"

	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
)

type todoHandler struct {
	db *bun.DB
}

func (h todoHandler) CreateTodo(c *fiber.Ctx) error {
	todoForm := dtos.CreateTodoForm{}
	err := c.BodyParser(&todoForm)
	// 1. validate request body
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// if len(todoForm.Text) == 0 {
	// 	return fiber.NewError(fiber.StatusUnprocessableEntity, "text is required")
	// }

	errors := validator.ValidateStruct(&todoForm)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalidate JSON",
			"errors":  errors,
		})
	}

	// 2. insert to database
	// todo := models.Todo{Title: todoForm.Text}
	todo := mapper.CreateTodoFormToModel(&todoForm)
	_, err = h.db.NewInsert().Model(todo).Exec(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	serialized := mapper.TodoToDto(todo)

	// 3. response result
	return c.Status(201).JSON(serialized)
}

func (h todoHandler) ListTodo(c *fiber.Ctx) error {
	filter := struct {
		// เป็น pointer เพื่อไม่ส่งค่ามาจะให้เป็น nil เพื่อแสดงข้อมูลทั้งหมด
		Completed *bool `query:"completed"`
	}{}
	err := c.QueryParser(&filter)

	// 1. validate request body
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// 2. select from database
	todos := []models.Todo{}
	q := h.db.NewSelect().Model(&todos).Order("created_at")
	if filter.Completed != nil {
		q.Where("is_done = ?", *filter.Completed)
	}
	err = q.Scan(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	serialized := mapper.TodosToDto(todos)

	// 3. response result
	return c.JSON(serialized)
}

func (h todoHandler) GetTodo(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	// 1. validate request body
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// 2. select from database by id
	todo := models.Todo{ID: id}
	err = h.db.NewSelect().Model(&todo).WherePK().Scan(c.UserContext())
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "todo with given id not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	serialized := mapper.TodoToDto(&todo)

	// 3. response result
	return c.JSON(serialized)
}

func (h todoHandler) UpdateTodoStatus(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	// 1. validate request body
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	updateTodo := dtos.UpdateTodoForm{}
	err = c.BodyParser(&updateTodo)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// 2. update to database by id
	// todo := models.Todo{ID: id, IsDone: updateTodo.Completed}
	todo := models.Todo{ID: id, IsDone: updateTodo.Completed}
	res, err := h.db.NewUpdate().
		Model(&todo).
		Column("is_done").
		SetColumn("updated_at", "DEFAULT").
		WherePK().
		Returning("*").
		Exec(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "todo with given id not found",
		})
	}

	serialized := mapper.TodoToDto(&todo)

	// 3. response result
	return c.JSON(serialized)
}

func (h todoHandler) DeleteTodo(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	// 1. validate request body
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// 2. delete from to database by id
	res, err := h.db.NewDelete().
		Model((*models.Todo)(nil)).
		Where("id = ?", id).
		Exec(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "todo with given id not found",
		})
	}

	// 3. response result
	return c.SendStatus(204)
}
