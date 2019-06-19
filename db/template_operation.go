package db

import (
	"github.com/hetianyi/gox"
	"github.com/hetianyi/gox/logger"
	"github.com/jinzhu/gorm"
)

// Query does a db query.
func Query(result interface{}, sql string, args ...interface{}) {

}

func Transaction(work func(tx *gorm.DB) error) error {
	db := GetDB()
	tx := db.Begin()
	gox.Try(func() {
		err := work(tx)
		if err != nil {
			if r := tx.Rollback(); r.Error != nil {
				logger.Error("roll back transaction failed: ", r.Error)
			}
			return
		}
		if r := tx.Commit(); r.Error != nil {
			logger.Error("commit transaction failed: ", r.Error)
		}
	}, func(e interface{}) {
		logger.Error("transaction error: ", e)
		if r := tx.Rollback(); r.Error != nil {
			logger.Error("roll back transaction failed: ", r.Error)
		}
	})
	return nil
}
