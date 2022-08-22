package ports

import (
	"goapi/pkg/common"
	"goapi/pkg/module/user/core/dto"
	"goapi/pkg/module/user/core/model"
)

// interface สำหรับ output port
type UserRepository interface {
	Create(m *model.User) error
	Find(page common.PagingRequest) (model.Users, *common.PagingResult, error)
	FindById(id string) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	UpdatePasswordById(id string, m *model.User) error
	DeleteById(id string) error
}

// interface สำหรับ input port
type UserService interface {
	Create(newUser dto.NewUserForm, reqId string) (*dto.UserResponse, error)
	List(page common.PagingRequest, reqId string) (dto.UserResponses, *common.PagingResult, error)
	Get(id string, reqId string) (*dto.UserResponse, error)
	UpdatePassword(id string, updateUser dto.UpdateUserPasswordForm, reqId string) (*dto.UserResponse, error)
	Delete(id string, reqId string) error
}
