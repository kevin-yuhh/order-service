package model

import (
	"time"

	"github.com/TRON-US/soter-order-service/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type Database struct {
	DB *xorm.EngineGroup
}

// New database connection.
func NewDatabase(config *config.Configuration) (*Database, error) {
	connections := []string{
		config.Database.ConnectionUriMaster,
		config.Database.ConnectionUriSlave,
	}

	group, err := xorm.NewEngineGroup("mysql", connections, xorm.RandomPolicy())
	if err != nil {
		return nil, err
	}

	group.SetMaxIdleConns(config.Database.MaxIdleConn)
	group.SetMaxOpenConns(config.Database.MaxOpenConn)
	group.SetConnMaxLifetime(time.Second * time.Duration(config.Database.MaxLifetime))

	return &Database{DB: group}, nil
}

// Close database connection.
func (db *Database) Close() (err error) {
	err = db.DB.Close()
	return
}
