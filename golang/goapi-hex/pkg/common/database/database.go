package database

import (
	"goapi-hex/pkg/core/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(model.Todo{})
	return db, nil
}

func CloseDB(db *gorm.DB) {
	sqlDB, _ := db.DB()
	sqlDB.Close()
}
