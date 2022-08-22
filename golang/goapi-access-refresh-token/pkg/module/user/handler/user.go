package handler

import (
	"goapi/pkg/common"
	"goapi/pkg/module/user/core/dto"
	"goapi/pkg/module/user/core/ports"
)

type UserHandler struct {
	serv ports.UserService
}

func NewUserHandler(serv ports.UserService) *UserHandler {
	return &UserHandler{serv}
}

// @Summary Add a new user
// @Description Add a new user
// @Tags User
// @Accept  json
// @Produce  json
// @Param user body swagger.CreateUserFrom true "User Data"
// @Failure 422 {object} swagdto.Error422{error=swagger.ErrCreateSampleData}
// @Failure 500 {object} swagdto.Error500
// @Success 201 {object} swagdto.Response{data=swagger.UserSampleData}
// @Router /users [post]
func (h UserHandler) CreateUser(c common.HContext) error {
	// แปลง JSON เป็น struct
	form := new(dto.NewUserForm)
	if err := c.BodyParser(form); err != nil {
		return common.ResponseError(c, common.ErrBodyParser)
	}
	// ส่งต่อไปให้ service ทำงาน
	user, err := h.serv.Create(*form, c.RequestId())
	if err != nil {
		// error จะถูกจัดการมาจาก service แล้ว
		return common.ResponseError(c, err)
	}

	// คืนค่า user ที่เพิ่งบันทึกเสร็จกลับไปในรูปแบบ JSON
	return common.ResponseCreated(c, "user", user)
}

// @Summary List all existing users
// @Description You can filter all existing users by listing them.
// @Tags User
// @Accept  json
// @Produce  json
// @Param page query int false "Go to a specific page number. Start with 1"
// @Param limit query int false "Page size for the data"
// @Param order query string false "Page order. Eg: text desc,createdAt desc"
// @Failure 400 {object} swagdto.Error400
// @Failure 500 {object} swagdto.Error500
// @Success 200 {object} swagdto.ResponseWithPage{data=swagger.UserSampleListData}
// @Router /users [get]
func (h UserHandler) ListUser(c common.HContext) error {
	page := common.Paginator(c)

	users, paging, err := h.serv.List(page, c.RequestId())

	if err != nil {
		return common.ResponseError(c, err)
	}

	return common.ResponsePage(c, "users", users, paging)
}

// @Summary Get a user
// @Description Get a specific user by id
// @Produce json
// @Tags User
// @Param id path string true "User ID"
// @Failure 400 {object} swagdto.Error400
// @Failure 404 {object} swagdto.Error404
// @Failure 500 {object} swagdto.Error500
// @Success 200 {object} swagdto.Response{data=swagger.UserSampleData}
// @Router /users/{id} [get]
func (h UserHandler) GetUser(c common.HContext) error {
	id := c.Param("id")

	user, err := h.serv.Get(id, c.RequestId())

	if err != nil {
		return common.ResponseError(c, err)
	}

	return common.ResponseOk(c, "user", user)
}

// @Summary Update a user password
// @Description Update a specific user password by id
// @Produce json
// @Tags User
// @Param id path string true "User ID"
// @Param user body swagger.UpdateUserPasswordForm true "User Password"
// @Failure 400 {object} swagdto.Error400
// @Failure 404 {object} swagdto.Error404
// @Failure 422 {object} swagdto.Error422{error=swagger.ErrUpdateSampleData}
// @Failure 500 {object} swagdto.Error500
// @Success 200 {object} swagdto.Response{data=swagger.UserSampleData}
// @Router /users/{id} [patch]
func (h UserHandler) UpdateUserPassword(c common.HContext) error {
	id := c.Param("id")

	form := dto.UpdateUserPasswordForm{}

	if err := c.BodyParser(&form); err != nil {
		return common.ResponseError(c, err)
	}

	user, err := h.serv.UpdatePassword(id, form, c.RequestId())

	if err != nil {
		return common.ResponseError(c, err)
	}

	return common.ResponseOk(c, "user", user)
}

// @Summary Delete a user
// @Description Delete a specific user by id
// @Produce  json
// @Tags User
// @Param id path string true "User ID"
// @Failure 400 {object} swagdto.Error400
// @Failure 404 {object} swagdto.Error404
// @Failure 500 {object} swagdto.Error500
// @Success 204
// @Router /users/{id} [delete]
func (h UserHandler) DeleteUser(c common.HContext) error {
	id := c.Param("id")

	err := h.serv.Delete(id, c.RequestId())

	if err != nil {
		return common.ResponseError(c, err)
	}

	return common.ResponseNoContent(c)
}
