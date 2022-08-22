package database

import (
	"fmt"
	"goapi/pkg/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GormDB struct {
	*gorm.DB
}

func NewGormDB(conf *config.Config) (*GormDB, error) {
	// Build a DSN e.g. postgres://username:password@host:port/dbName
	dsn := fmt.Sprintf("%v://%v:%v@%v:%v/%v?sslmode=%v",
		conf.Db.Driver,
		conf.Db.Username,
		conf.Db.Password,
		conf.Db.Host,
		conf.Db.Port,
		conf.Db.Database,
		conf.Db.Sslmode)

	gcnf := &gorm.Config{}

	if conf.App.Mode == "production" {
		gcnf.Logger = logger.Default.LogMode(logger.Silent)
	} else {
		gcnf.Logger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(postgres.Open(dsn), gcnf)
	if err != nil {
		return nil, err
	}
	return &GormDB{db}, nil
}

func (g *GormDB) CloseDB() error {
	sqlDB, err := g.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
