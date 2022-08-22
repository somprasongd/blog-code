package model

import "database/sql/driver"

type UserRole string

const (
	ADMIN UserRole = "admin"
	USER  UserRole = "user"
)

func (e *UserRole) Scan(value interface{}) error {
	*e = UserRole(value.(string))
	return nil
}

func (e UserRole) Value() (driver.Value, error) {
	return string(e), nil
}

func (e UserRole) String() string {
	switch e {
	case ADMIN:
		return "admin"
	default:
		return "user"
	}
}
