package models

import "time"

type Todo struct {
	ID     int    `bun:"id,pk,autoincrement"`
	Title  string `bun:"title,nullzero,notnull"`
	IsDone bool   `bun:"is_done,notnull,default:false"`

	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	DeletedAt time.Time `bun:",soft_delete,nullzero"`
}
