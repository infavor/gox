package main

import (
	"fmt"
	"github.com/hetianyi/gox/db"
	"github.com/hetianyi/gox/logger"
	"github.com/hetianyi/gox/uuid"
	"github.com/jinzhu/gorm"
	"time"
)

func init() {
	logger.Init(nil)
}

type Article1 struct {
	Id    int    `gorm:"column:id"`
	Title string `gorm:"column:ma_title"`
}

func main() {
	driverName := "mysql"
	connectString := "root:123456@tcp(192.168.1.142:3306)/mao"
	if err := db.InitConnection(driverName, connectString); err != nil {
		logger.Fatal(err)
	}
	db.SetConfig(10, 50, time.Hour*12)

	err := db.Transaction(func(tx *gorm.DB) error {
		a := &Article1{
			Title: uuid.UUID(),
		}
		ret := tx.Table("article1").Create(a)
		if ret.Error != nil {
			return ret.Error
		}
		b := &Article1{
			Id:    a.Id,
			Title: uuid.UUID(),
		}
		ret = tx.Table("article1").Create(b)
		if ret.Error != nil {
			return ret.Error
		}
		if ret.RowsAffected == 0 {
			logger.Error("error insert record")
		}
		return nil
	})
	if err != nil {
		logger.Error(err)
	}
	fmt.Println()
}
