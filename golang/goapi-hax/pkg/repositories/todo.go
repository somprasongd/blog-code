package repositories

import (
	"errors"
	"goapi-hax/pkg/core/domain"
	"goapi-hax/pkg/core/ports"

	"gorm.io/gorm"
)

type todoRepositoryDB struct {
	db *gorm.DB
}

func NewTodoRepositoryDB(db *gorm.DB) ports.TodoRepository {
	return &todoRepositoryDB{
		db: db,
	}
}

func (r todoRepositoryDB) Create(t *domain.Todo) error {
	tx := r.db.Create(&t)
	if err := tx.Error; err != nil {
		return err
	}
	return nil
}

func (r todoRepositoryDB) Find(conditions map[string]interface{}) ([]domain.Todo, error) {
	todos := []domain.Todo{}
	// เพิ่ม .Where()
	tx := r.db.Where(conditions).Find(&todos)

	if err := tx.Error; err != nil {
		return nil, err
	}
	return todos, nil
}

func (r todoRepositoryDB) FindById(id int) (*domain.Todo, error) {
	todo := &domain.Todo{}
	tx := r.db.First(todo, id)
	if err := tx.Error; err != nil {
		// handle error not found
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("NOT_FOUND")
		}
		return nil, err
	}
	return todo, nil
}

func (r todoRepositoryDB) UpdateStatusById(id int, isDone bool) error {
	tx := r.db.Model(domain.Todo{ID: id}).Update("is_done", isDone)
	if err := tx.Error; err != nil {
		return err
	}
	// handle not found error
	if tx.RowsAffected == 0 {
		return errors.New("NOT_FOUND")
	}
	return nil
}

func (r todoRepositoryDB) DeleteById(id int) error {
	tx := r.db.Delete(&domain.Todo{}, id)
	if err := tx.Error; err != nil {
		return err
	}
	// handle not found error
	if tx.RowsAffected == 0 {
		return errors.New("NOT_FOUND")
	}
	return nil
}
