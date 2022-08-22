package model

type Todo struct {
	ID   uint `gorm:"primaryKey"`
	Text string
	Done bool
}
