package service

import (
	"errors"
	"goapi/pkg/common"
	"goapi/pkg/common/logger"
	"goapi/pkg/config"
	"goapi/pkg/module/auth/core/dto"
	"goapi/pkg/module/auth/core/model"
	"goapi/pkg/module/auth/core/ports"
	"goapi/pkg/util"
)

var (
	ErrUserEmailDuplication = common.NewBadRequestError("email already exists")
	ErrHashPassword         = common.NewUnexpectedError("error occurred while hashing password")
	ErrLogin                = common.NewUnauthorizedError("the email or password are incorrect")
	ErrGenerateToken        = common.NewUnexpectedError("error occurred while generating token")
	ErrValidateToken        = common.NewUnexpectedError("error occurred while validating token")
	ErrNoToken              = common.NewUnauthorizedError("the token is required")
	ErrInvalidToken         = common.NewUnauthorizedError("the token is invalid")
	ErrUserNotfound         = common.NewUnauthorizedError("user not found")
	ErrUserPasswordNotMatch = common.NewBadRequestError("password is not macth")
)

type authService struct {
	config *config.Config
	repo   ports.AuthRepository
}

func NewAuthService(config *config.Config, repo ports.AuthRepository) ports.AuthService {
	return &authService{config, repo}
}

func (s authService) Register(form dto.RegisterForm, reqId string) error {
	// validate
	if err := common.ValidateDto(form); err != nil {
		return common.NewInvalidError(err.Error())
	}

	u, err := s.repo.FindUserByEmail(form.Email)

	if err != nil && !errors.Is(err, common.ErrRecordNotFound) {
		logger.ErrorWithReqId(err.Error(), reqId)
		return common.ErrDbQuery
	}

	if u != nil {
		return ErrUserEmailDuplication
	}

	auth := model.User{Email: form.Email}
	hashPwd, err := util.HashPassword(form.Password)
	if err != nil {
		logger.ErrorWithReqId(err.Error(), reqId)
		return ErrHashPassword
	}
	auth.Password = hashPwd

	err = s.repo.CreateUser(&auth)
	if err != nil {
		logger.ErrorWithReqId(err.Error(), reqId)
		return common.ErrDbInsert
	}

	return nil
}

func (s authService) Login(form dto.LoginForm, reqId string) (*dto.AuthResponse, error) {
	// validate form
	err := common.ValidateDto(form)
	if err != nil {
		return nil, common.NewInvalidError(err.Error())
	}
	// ค้นหาจาก email
	user, err := s.repo.FindUserByEmail(form.Email)
	if err != nil {
		if errors.Is(err, common.ErrRecordNotFound) {
			return nil, ErrLogin
		}
		logger.ErrorWithReqId(err.Error(), reqId)
		return nil, common.ErrDbQuery
	}
	// ตรวจสอบรหัสผ่าน ตรงกันหรือไม่
	match := util.CheckPasswordHash(form.Password, user.Password)

	if !match {
		return nil, ErrLogin
	}
	// สร้าง jwt token
	token, err := util.GenerateToken(user.ID.String(), user.Email, user.Role.String(), s.config.Token.SecretKey)

	if err != nil {
		logger.ErrorWithReqId(err.Error(), reqId)
		return nil, ErrGenerateToken
	}
	// ตอบกลับไปพร้อมข้อมูล user
	serialized := dto.AuthResponse{
		User: dto.UserInfo{
			ID:    user.ID.String(),
			Email: user.Email,
			Role:  user.Role.String(),
		},
		Token: token,
	}
	return &serialized, nil
}

func (s authService) Profile(email string, reqId string) (*dto.UserInfo, error) {
	// validate
	if email == "" {
		return nil, ErrUserNotfound
	}

	user, err := s.repo.FindUserByEmail(email)

	if err != nil {
		if errors.Is(err, common.ErrRecordNotFound) {
			return nil, ErrUserNotfound
		}
		logger.ErrorWithReqId(err.Error(), reqId)
		return nil, common.ErrDbQuery
	}

	serialized := dto.UserInfo{
		ID:    user.ID.String(),
		Email: user.Email,
		Role:  user.Role.String(),
	}

	return &serialized, nil
}

func (s authService) UpdateProfile(email string, form dto.UpdateProfileForm, reqId string) (*dto.UserInfo, error) {
	// validate
	err := common.ValidateDto(form)
	if err != nil {
		return nil, common.NewInvalidError(err.Error())
	}

	user, err := s.repo.FindUserByEmail(email)
	if err != nil {
		if errors.Is(err, common.ErrRecordNotFound) {
			return nil, ErrUserNotfound
		}
		logger.ErrorWithReqId(err.Error(), reqId)
		return nil, common.ErrDbQuery
	}

	match := util.CheckPasswordHash(form.PasswordOld, user.Password)

	if !match {
		return nil, ErrUserPasswordNotMatch
	}

	hashPwd, err := util.HashPassword(form.PasswordNew)

	if err != nil {
		logger.ErrorWithReqId(err.Error(), reqId)
		return nil, ErrHashPassword
	}

	user.Password = hashPwd

	err = s.repo.SaveProfile(user)
	if err != nil {
		logger.ErrorWithReqId(err.Error(), reqId)
		return nil, common.ErrDbUpdate
	}

	serialized := dto.UserInfo{
		ID:    user.ID.String(),
		Email: user.Email,
		Role:  user.Role.String(),
	}

	return &serialized, nil
}
