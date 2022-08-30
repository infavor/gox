package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/infavor/gox/logger"
	"github.com/jinzhu/gorm"
	"time"
)

var (
	_driverName string
	_args       []interface{}
	managedDB   *gorm.DB
)

// InitConnection initializes db connection.
func InitConnection(dialect string, args ...interface{}) error {
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
	managedDB.LogMode(logger.GetLogLevel() > logger.DebugLevel)
	managedDB.SetLogger()
	return nil
}

// SetConfig sets configuration of connection pool.
func SetConfig(maxIdleConns int, maxOpenConns int, connMaxLifetime time.Duration) {
	managedDB.DB().SetMaxIdleConns(maxIdleConns)
	managedDB.DB().SetMaxOpenConns(maxOpenConns)
	managedDB.DB().SetConnMaxLifetime(connMaxLifetime)
}

// GetDB returns *gorm.DB for db operation.
func GetDB() *gorm.DB {
	return managedDB
}

// Reconnect tries to reconnect to db.
func Reconnect() error {
	if managedDB != nil {
		managedDB.DB().Close()
	}
	managedDB = nil
	return InitConnection(_driverName, _args...)
}
