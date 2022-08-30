package db

import (
	"github.com/infavor/gox"
	"github.com/infavor/gox/logger"
	"github.com/jinzhu/gorm"
)

// Query does a db query.
func Query(work func(sql *gorm.DB) error) error {
	return work(GetDB())
}

// Transaction does a db transaction.
func Transaction(work func(tx *gorm.DB) error) error {
	tx := GetDB().Begin()
	var err error
	gox.Try(func() {
		err = work(tx)
		if err != nil {
			if r := tx.Rollback(); r.Error != nil {
				logger.Error("roll back transaction failed: ", r.Error)
			}
			return
		}
		if r := tx.Commit(); r.Error != nil {
			logger.Error("commit transaction failed: ", r.Error)
			err = r.Error
		}
	}, func(e interface{}) {
		logger.Error("transaction error: ", e)
		if r := tx.Rollback(); r.Error != nil {
			logger.Error("roll back transaction failed: ", r.Error)
		}
		err = e.(error)
	})
	return err
}
