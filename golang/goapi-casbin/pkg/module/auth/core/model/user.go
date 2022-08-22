package model

import (
	"time"

	"github.com/gofrs/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Email     string
	Password  string
	Role      UserRole `sql:"user_role" gorm:"default:'user'"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *User) Promote() {
	u.Role = ADMIN
}

func (u *User) Demote() {
	u.Role = USER
}

type Users []*User
