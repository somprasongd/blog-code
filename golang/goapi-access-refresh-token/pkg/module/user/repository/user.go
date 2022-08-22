package repository

import (
	"errors"
	"fmt"
	"goapi/pkg/common"
	"goapi/pkg/module/user/core/model"
	"goapi/pkg/module/user/core/ports"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userRepositoryDB struct {
	db *gorm.DB
}

func NewUserRepositoryDB(db *gorm.DB) ports.UserRepository {
	return &userRepositoryDB{db}
}

func (r userRepositoryDB) Create(user *model.User) error {
	return r.db.Create(&user).Error
}

func (r userRepositoryDB) Find(page common.PagingRequest) (model.Users, *common.PagingResult, error) {
	users := []*model.User{}
	db := r.db

	pg := common.Pagination{
		PagingRequest: page,
		Query:         db,
		Records:       &users,
	}
	paging, err := pg.Paginate()

	if err != nil {
		return nil, nil, err
	}
	fmt.Println(users)
	return users, paging, nil
}

func (r userRepositoryDB) FindById(id string) (*model.User, error) {
	user := model.User{}
	db := r.db.Where("id = ?", id).First(&user)
	if err := db.Error; err != nil {
		// handle error not found
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r userRepositoryDB) FindByEmail(email string) (*model.User, error) {
	user := model.User{}
	db := r.db.Where("email = ?", email).First(&user)
	if err := db.Error; err != nil {
		// handle error not found
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r userRepositoryDB) UpdatePasswordById(id string, user *model.User) error {
	db := r.db.Model(&user).
		Clauses(clause.Returning{}).
		Where("id = ?", id).
		Updates(&user)
	if err := db.Error; err != nil {
		return err
	}
	// handle not found error
	if db.RowsAffected == 0 {
		return common.ErrRecordNotFound
	}
	return nil
}

func (r userRepositoryDB) DeleteById(id string) error {
	db := r.db.Where("id = ?", id).Delete(&model.User{})
	if err := db.Error; err != nil {
		return err
	}
	// handle not found error
	if db.RowsAffected == 0 {
		return common.ErrRecordNotFound
	}
	return nil
}
