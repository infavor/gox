package db

import (
	"github.com/jinzhu/gorm"
	"time"
)

var (
	_driverName string
	_args       []interface{}
	managedDB   *gorm.DB
)

func Open(dialect string, args ...interface{}) error {
	_driverName = dialect
	_args = args
	db, err := gorm.Open(dialect, args...)
	if err != nil {
		return err
	}
	if err = db.DB().Ping(); err != nil {
		return err
	}
	managedDB = db
	return nil
}

func SetConfig(maxIdleConns int, maxOpenConns int, connMaxLifetime time.Duration) {
	managedDB.DB().SetMaxIdleConns(maxIdleConns)
	managedDB.DB().SetMaxOpenConns(maxOpenConns)
	managedDB.DB().SetConnMaxLifetime(connMaxLifetime)
}

func GetDB() *gorm.DB {
	return managedDB
}

func Reconnect() error {
	if managedDB != nil {
		managedDB.DB().Close()
	}
	managedDB = nil
	return Open(_driverName, _args...)
}
