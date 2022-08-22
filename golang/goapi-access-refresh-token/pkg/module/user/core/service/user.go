package service

import (
	"errors"
	"fmt"
	"goapi/pkg/common"
	"goapi/pkg/common/logger"
	"goapi/pkg/module/user/core/dto"
	"goapi/pkg/module/user/core/mapper"
	"goapi/pkg/module/user/core/ports"
	"goapi/pkg/util"
)

var (
	// ErrUserNotFoundById user not found error when find with id
	ErrUserNotFoundById     = common.NewNotFoundError("user with given id not found")
	ErrHashPassword         = common.NewUnexpectedError("hash password error")
	ErrUserEmailDuplication = common.NewBadRequestError("duplicate email")
)

type userService struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) ports.UserService {
	return &userService{repo}
}

func (s userService) Create(form dto.NewUserForm, reqId string) (*dto.UserResponse, error) {
	// validate
	if err := common.ValidateDto(form); err != nil {
		return nil, common.NewInvalidError(err.Error())
	}

	u, err := s.repo.FindByEmail(form.Email)

	if err != nil && !errors.Is(err, common.ErrRecordNotFound) {
		logger.ErrorWithReqId(err.Error(), reqId)
		return nil, common.ErrDbQuery
	}

	if u != nil {
		return nil, ErrUserEmailDuplication
	}

	user := mapper.CreateUserFormToModel(form)
	hashPwd, err := util.HashPassword(form.Password)
	if err != nil {
		logger.ErrorWithReqId(err.Error(), reqId)
		return nil, ErrHashPassword
	}
	fmt.Println(form.Password, hashPwd)
	user.Password = hashPwd

	err = s.repo.Create(user)
	if err != nil {
		logger.ErrorWithReqId(err.Error(), reqId)
		return nil, common.ErrDbInsert
	}

	serialized := mapper.UserToDto(user)

	return serialized, nil
}

func (s userService) List(page common.PagingRequest, reqId string) (dto.UserResponses, *common.PagingResult, error) {
	users, paging, err := s.repo.Find(page)
	if err != nil {
		logger.ErrorWithReqId(err.Error(), reqId)
		return nil, nil, common.ErrDbQuery
	}

	serialized := mapper.UsersToDto(users)
	return serialized, paging, nil
}

func (s userService) Get(id string, reqId string) (*dto.UserResponse, error) {
	// validate id format
	if err := s.validateId(id); err != nil {
		return nil, common.ErrIdFormat
	}

	user, err := s.repo.FindById(id)
	if err != nil {
		if errors.Is(err, common.ErrRecordNotFound) {
			return nil, ErrUserNotFoundById
		}
		logger.ErrorWithReqId(err.Error(), reqId)
		return nil, common.ErrDbQuery
	}

	serialized := mapper.UserToDto(user)
	return serialized, nil
}

func (s userService) UpdatePassword(id string, form dto.UpdateUserPasswordForm, reqId string) (*dto.UserResponse, error) {
	// validate id format
	if err := s.validateId(id); err != nil {
		return nil, common.ErrIdFormat
	}

	err := common.ValidateDto(form)
	if err != nil {
		return nil, common.NewInvalidError(err.Error())
	}

	user, err := s.repo.FindById(id)
	if err != nil {
		if errors.Is(err, common.ErrRecordNotFound) {
			return nil, ErrUserNotFoundById
		}
		logger.ErrorWithReqId(err.Error(), reqId)
		return nil, common.ErrDbQuery
	}

	hashPwd, err := util.HashPassword(form.Password)

	if err != nil {
		logger.ErrorWithReqId(err.Error(), reqId)
		return nil, ErrHashPassword
	}

	user.Password = hashPwd

	err = s.repo.UpdatePasswordById(id, user)
	if err != nil {
		if errors.Is(err, common.ErrRecordNotFound) {
			return nil, ErrUserNotFoundById
		}
		logger.ErrorWithReqId(err.Error(), reqId)
		return nil, common.ErrDbUpdate
	}

	serialized := mapper.UserToDto(user)
	return serialized, nil
}

func (s userService) Delete(id string, reqId string) error {
	// validate id format
	if err := s.validateId(id); err != nil {
		return common.ErrIdFormat
	}

	err := s.repo.DeleteById(id)
	if err != nil {
		if errors.Is(err, common.ErrRecordNotFound) {
			return ErrUserNotFoundById
		}
		logger.ErrorWithReqId(err.Error(), reqId)
		return common.ErrDbDelete
	}

	return nil
}

func (s userService) validateId(id string) error {
	return common.ValidateDto(&dto.UserId{ID: id})
}
