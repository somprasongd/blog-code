package model

import (
	"time"

	"github.com/gofrs/uuid"
)

type Todo struct {
	ID        uuid.UUID  `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Text      string     `gorm:"not null"`
	Status    TodoStatus `sql:"todo_status" gorm:"default:'open'"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (t *Todo) Open() {
	t.Status = OPEN
}

func (t *Todo) Done() {
	t.Status = DONE
}

type Todos []*Todo
