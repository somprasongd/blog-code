package model

import "database/sql/driver"

type TodoStatus string

const (
	OPEN TodoStatus = "open"
	// IN_PROGRESS TodoStatus = "in_progress"
	DONE TodoStatus = "done"
)

func (e *TodoStatus) Scan(value interface{}) error {
	*e = TodoStatus(value.(string))
	return nil
}

func (e TodoStatus) Value() (driver.Value, error) {
	return string(e), nil
}

func (e TodoStatus) String() string {
	switch e {
	// case IN_PROGRESS:
	// 	return "in_progress"
	case DONE:
		return "done"
	default:
		return "open"
	}
}
