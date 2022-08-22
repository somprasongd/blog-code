package handler

import (
	"goapi/pkg/common"
	"goapi/pkg/module/auth/core/dto"
	"goapi/pkg/module/auth/core/ports"

	"github.com/golang-jwt/jwt/v4"
)

type AuthHandler interface {
	Register(common.HContext) error
	Login(c common.HContext) error
	Profile(c common.HContext) error
	UpdateProfile(c common.HContext) error
	RefreshToken(c common.HContext) error
	RevokeToken(c common.HContext) error
}

type authHandler struct {
	serv ports.AuthService
}

func NewAuthHandler(serv ports.AuthService) AuthHandler {
	return &authHandler{serv}
}

// @Summary Register a new user
// @Description Register a new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body swagger.RegisterForm true "User Data"
// @Failure 422 {object} swagdto.Error422{error=swagger.ErrRegisterSampleData}
// @Failure 500 {object} swagdto.Error500
// @Success 201
// @Router /auth/register [post]
func (h authHandler) Register(c common.HContext) error {
	// แปลง JSON เป็น struct
	form := new(dto.RegisterForm)
	if err := c.BodyParser(form); err != nil {
		return common.ResponseError(c, common.ErrBodyParser)
	}
	// ส่งต่อไปให้ service ทำงาน
	err := h.serv.Register(*form, c.RequestId())
	if err != nil {
		// error จะถูกจัดการมาจาก service แล้ว
		return common.ResponseError(c, err)
	}

	return common.ResponseCreated(c, "", nil)
}

// @Summary Login
// @Description Login
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param user body swagger.LoginForm true "Login Data"
// @Failure 401 {object} swagdto.Error401
// @Failure 422 {object} swagdto.Error422{error=swagger.ErrLoginSampleData}
// @Failure 500 {object} swagdto.Error500
// @Success 200 {object} swagdto.Response{data=swagger.AuthSampleData}
// @Router /auth/login [post]
func (h authHandler) Login(c common.HContext) error {
	// แปลง JSON เป็น struct
	form := new(dto.LoginForm)
	if err := c.BodyParser(form); err != nil {
		return common.ResponseError(c, common.ErrBodyParser)
	}
	// ส่งต่อไปให้ service ทำงาน
	auth, err := h.serv.Login(*form, c.RequestId())
	if err != nil {
		// error จะถูกจัดการมาจาก service แล้ว
		return common.ResponseError(c, err)
	}

	return common.ResponseOk(c, "auth", auth)
}

// @Summary Get a user profile
// @Description Get a specific user by id
// @Produce json
// @Tags Auth
// @Param Authorization header string true "Bearer"
// @Failure 401 {object} swagdto.Error401
// @Failure 500 {object} swagdto.Error500
// @Success 200 {object} swagdto.Response{data=swagger.UserSampleData}
// @Router /auth/profile [get]
func (h authHandler) Profile(c common.HContext) error {
	u := c.Locals("user").(jwt.MapClaims)
	email := u["email"].(string)

	user, err := h.serv.Profile(email, c.RequestId())

	if err != nil {
		return common.ResponseError(c, err)
	}

	return common.ResponseOk(c, "user", user)
}

// @Summary Update a user password
// @Description Update a user password
// @Produce json
// @Tags User
// @Param Authorization header string true "Bearer"
// @Param user body swagger.UpdateProfileForm true "User Password"
// @Failure 400 {object} swagdto.Error400
// @Failure 404 {object} swagdto.Error404
// @Failure 422 {object} swagdto.Error422{error=swagger.ErrUpdateSampleData}
// @Failure 500 {object} swagdto.Error500
// @Success 200 {object} swagdto.Response{data=swagger.UserSampleData}
// @Router /users/{id} [patch]
func (h authHandler) UpdateProfile(c common.HContext) error {
	u := c.Locals("user").(jwt.MapClaims)
	email := u["email"].(string)

	form := dto.UpdateProfileForm{}

	if err := c.BodyParser(&form); err != nil {
		return common.ResponseError(c, err)
	}

	user, err := h.serv.UpdateProfile(email, form, c.RequestId())

	if err != nil {
		return common.ResponseError(c, err)
	}

	return common.ResponseOk(c, "user", user)
}

// @Summary Refresh Token
// @Description Generate new access and refresh token from refresh token
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param user body swagger.RefreshForm true "Refresh Token Data"
// @Failure 401 {object} swagdto.Error401
// @Failure 422 {object} swagdto.Error422{error=swagger.ErrLoginSampleData}
// @Failure 500 {object} swagdto.Error500
// @Success 200 {object} swagdto.Response{data=swagger.AuthSampleData}
// @Router /auth/refresh [post]
func (h authHandler) RefreshToken(c common.HContext) error {
	// แปลง JSON เป็น struct
	form := new(dto.RefreshForm)
	if err := c.BodyParser(form); err != nil {
		return common.ResponseError(c, common.ErrBodyParser)
	}
	// ส่งต่อไปให้ service ทำงาน
	auth, err := h.serv.RefreshToken(*form, c.RequestId())
	if err != nil {
		// error จะถูกจัดการมาจาก service แล้ว
		return common.ResponseError(c, err)
	}

	return common.ResponseOk(c, "auth", auth)
}

// @Summary Revoke Token
// @Description Remove token id in redis
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param user body swagger.RefreshForm true "Refresh Token Data"
// @Failure 401 {object} swagdto.Error401
// @Failure 422 {object} swagdto.Error422{error=swagger.ErrLoginSampleData}
// @Failure 500 {object} swagdto.Error500
// @Success 204
// @Router /auth/revoke [post]
func (h authHandler) RevokeToken(c common.HContext) error {
	// แปลง JSON เป็น struct
	form := new(dto.RefreshForm)
	if err := c.BodyParser(form); err != nil {
		return common.ResponseError(c, common.ErrBodyParser)
	}
	// ส่งต่อไปให้ service ทำงาน
	err := h.serv.RevokeToken(*form, c.RequestId())
	if err != nil {
		// error จะถูกจัดการมาจาก service แล้ว
		return common.ResponseError(c, err)
	}

	return common.ResponseNoContent(c)
}
