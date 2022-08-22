package database

import (
	"database/sql"
	"fmt"
	"goapi/pkg/config"

	_ "github.com/lib/pq"
)

type SqlDB struct {
	*sql.DB
}

func NewDB(conf *config.Config) (*SqlDB, error) {
	// Build a DSN e.g. postgres://username:password@host:port/dbName
	dsn := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=%v",
		conf.Db.Username,
		conf.Db.Password,
		conf.Db.Host,
		conf.Db.Port,
		conf.Db.Database,
		conf.Db.Sslmode)

	db, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, err
	}

	return &SqlDB{db}, nil
}

func (db *SqlDB) CloseDB() error {
	return db.Close()
}
