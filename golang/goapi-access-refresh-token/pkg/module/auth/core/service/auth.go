package service

import (
	"encoding/json"
	"errors"
	"goapi/pkg/common"
	"goapi/pkg/common/logger"
	"goapi/pkg/config"
	"goapi/pkg/module/auth/core/dto"
	"goapi/pkg/module/auth/core/ports"
	"goapi/pkg/module/user/core/model"
	"goapi/pkg/util"

	"github.com/gofrs/uuid"
)

var (
	ErrUserEmailDuplication = common.NewBadRequestError("email already exists")
	ErrHashPassword         = common.NewUnexpectedError("error occurred while hashing password")
	ErrLogin                = common.NewUnauthorizedError("the email or password are incorrect")
	ErrGenerateAccessToken  = common.NewUnexpectedError("error occurred while generating token")
	ErrGenerateRefreshToken = common.NewUnexpectedError("error occurred while generating refresh token")
	ErrValidateToken        = common.NewUnexpectedError("error occurred while validating token")
	ErrNoToken              = common.NewUnauthorizedError("the token is required")
	ErrInvalidToken         = common.NewUnauthorizedError("the token is invalid")
	ErrUserNotfound         = common.NewUnauthorizedError("user not found")
	ErrUserPasswordNotMatch = common.NewBadRequestError("password is not macth")
	ErrInvalidRefreshToken  = common.NewUnauthorizedError("the refresh token is invalid or expired")
	ErrUnmarshalPayload     = common.NewUnexpectedError("error occurred while convert user data")
)

type authService struct {
	config    *config.Config
	repo      ports.AuthRepository
	tokenRepo ports.TokenRepository
}

func NewAuthService(config *config.Config, repo ports.AuthRepository, tokenRepo ports.TokenRepository) ports.AuthService {
	return &authService{config, repo, tokenRepo}
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
	hashPwd, err := util.HashPasswordArgon(form.Password)
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
	// ???????????????????????? email
	user, err := s.repo.FindUserByEmail(form.Email)
	if err != nil {
		if errors.Is(err, common.ErrRecordNotFound) {
			return nil, ErrLogin
		}
		return nil, common.ErrDbQuery
	}
	// ????????????????????????????????????????????? ???????????????????????????????????????
	match := util.CheckPasswordHashArgon(form.Password, user.Password)

	if !match {
		return nil, ErrLogin
	}

	tokenId, _ := uuid.NewV4()
	payload := map[string]any{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"role":    user.Role.String(),
	}

	// ??????????????? refresh token
	refreshToken, refreshExpiresAt, err := util.GenerateToken(tokenId.String(), nil, s.config.Token.RefreshSecretKey, s.config.Token.RefreshExpires)

	if err != nil {
		logger.ErrorWithReqId(err.Error(), reqId)
		return nil, ErrGenerateRefreshToken
	}

	// ?????????????????? user ?????? redis
	err = s.tokenRepo.SetToken(tokenId.String(), payload, s.config.Token.RefreshExpires)
	if err != nil {
		logger.ErrorWithReqId(err.Error(), reqId)
		return nil, ErrGenerateRefreshToken
	}

	// ??????????????? access token
	accessToken, accessExpiresAt, err := util.GenerateToken(tokenId.String(), payload, s.config.Token.AccessSecretKey, s.config.Token.AccessExpires)

	if err != nil {
		logger.ErrorWithReqId(err.Error(), reqId)
		return nil, ErrGenerateAccessToken
	}

	// ???????????????????????????????????????????????????????????? user
	serialized := dto.AuthResponse{
		User: &dto.UserInfo{
			ID:    user.ID.String(),
			Email: user.Email,
			Role:  user.Role.String(),
		},
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessExpiresAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshExpiresAt,
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

func (s authService) RefreshToken(form dto.RefreshForm, reqId string) (*dto.AuthResponse, error) {
	// validate form
	err := common.ValidateDto(form)
	if err != nil {
		return nil, common.NewInvalidError(err.Error())
	}

	// ????????????????????? refresh token ?????????????????? valid ?????????????????????
	cliams, err := util.ValidateToken(form.Token, s.config.Token.RefreshSecretKey)

	if err != nil {
		logger.ErrorWithReqId(err.Error(), reqId)
		return nil, ErrInvalidRefreshToken
	}

	// ????????? token id ?????????????????? redis
	tokenId := cliams["sub"].(string)
	encodedUser, err := s.tokenRepo.GetToken(tokenId)

	if err != nil {
		logger.ErrorWithReqId(err.Error(), reqId)
		return nil, ErrInvalidRefreshToken
	}

	// ???????????????????????????????????????????????????????????????????????? ??????????????????????????????????????????????????????
	if encodedUser == "" {
		return nil, ErrInvalidRefreshToken
	}

	// ????????????????????? user ?????????????????????????????? payload
	payload := map[string]any{}
	err = json.Unmarshal([]byte(encodedUser), &payload)
	if err != nil {
		logger.ErrorWithReqId(err.Error(), reqId)
		return nil, ErrUnmarshalPayload
	}

	// ??????????????? tokenId ????????????
	newTkId, _ := uuid.NewV4()

	// ??????????????? refresh token ????????????
	refreshToken, refreshExpiresAt, err := util.GenerateToken(newTkId.String(), nil, s.config.Token.RefreshSecretKey, s.config.Token.RefreshExpires)

	if err != nil {
		logger.ErrorWithReqId(err.Error(), reqId)
		return nil, ErrGenerateRefreshToken
	}

	// ?????????????????? user ?????? redis
	err = s.tokenRepo.SetToken(newTkId.String(), payload, s.config.Token.RefreshExpires)
	if err != nil {
		logger.ErrorWithReqId(err.Error(), reqId)
		return nil, ErrGenerateRefreshToken
	}
	// ?????? tokenId ???????????? ???????????????????????????????????????
	s.tokenRepo.DeleteToken(tokenId)

	// ??????????????? access token ????????????
	accessToken, accessExpiresAt, err := util.GenerateToken(newTkId.String(), payload, s.config.Token.AccessSecretKey, s.config.Token.AccessExpires)

	if err != nil {
		logger.ErrorWithReqId(err.Error(), reqId)
		return nil, ErrGenerateAccessToken
	}

	// ????????? token ??????????????????????????????
	serialized := dto.AuthResponse{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessExpiresAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshExpiresAt,
	}
	return &serialized, nil
}

func (s authService) RevokeToken(form dto.RefreshForm, reqId string) error {
	// validate form
	err := common.ValidateDto(form)
	if err != nil {
		return common.NewInvalidError(err.Error())
	}

	// ????????????????????? refresh token ?????????????????? valid ?????????????????????
	cliams, err := util.ValidateToken(form.Token, s.config.Token.RefreshSecretKey)

	if err != nil {
		logger.ErrorWithReqId(err.Error(), reqId)
		return ErrInvalidRefreshToken
	}

	// ????????? token id ?????????????????? redis
	tokenId := cliams["sub"].(string)
	s.tokenRepo.DeleteToken(tokenId)

	return nil
}
