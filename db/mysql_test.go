package db_test

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hetianyi/gox/convert"
	"github.com/hetianyi/gox/db"
	"github.com/hetianyi/gox/logger"
	"github.com/hetianyi/gox/timer"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"testing"
	"time"
)

type Article struct {
	Id    int    `gorm:"column:id"`
	Title string `gorm:"column:ma_title"`
}

func init() {
	logger.Init(nil)
	db.RegisterSQL("query1", `select id, ma_title from article where id=?`)
}

var total = 0

func TestMySQL(t *testing.T) {
	driverName := "mysql"
	connectString := "root:123456@tcp(127.0.0.1:3306)/mao"
	if err := db.InitConnection(driverName, connectString); err != nil {
		logger.Fatal(err)
	}
	timer.Start(0, 0, time.Second, func(t *timer.Timer) {
		fmt.Println(convert.IntToStr(total) + "/s")
		total = 0
	})

	for true {
		err := db.Query(func(sql *gorm.DB) error {
			rows, err := sql.Raw(db.GetSQL("query1"), 6).Rows()
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				a := &Article{}
				err := sql.ScanRows(rows, a)
				if err != nil {
					return err
				}
				// logger.Info(fmt.Sprintf("id: %d, title: %s", a.Id, a.Title))
			}
			return nil
		})
		if err != nil {
			logger.Error(err)
		}
		total++
	}
}

func TestMysqlQuery(t *testing.T) {
	db, err := sql.Open("mysql", "root:123456@tcp(192.168.1.142:3306)/mao")
	if err != nil {
		logger.Fatal(err)
	}
	rows, err := db.Query("select id, ma_title from article a where a.id = 24")
	if err != nil {
		logger.Error(err)
		return
	}
	cols, _ := rows.Columns()
	logger.Info("columns: ", cols)
	for rows.Next() {
		var id int
		var title string
		if err = rows.Scan(&id, &title); err != nil {
			logger.Error(err)
			return
		}
		logger.Info("id=", id, " title=", title)
	}
}

func TestMysqlQuery1(t *testing.T) {
	db, err := gorm.Open("mysql", "root:123456@tcp(192.168.1.142:3306)/mao")
	if err != nil {
		logger.Fatal(err)
	}
	m := &Article{}
	ret := db.Table("article").Select("ma_title").Where("id = ?", 24).Scan(m)
	if ret.Error != nil {
		logger.Error(err)
		return
	}
	logger.Info(m)
}

func TestSqliteQuery(t *testing.T) {
	db, err := gorm.Open("sqlite3", "/home/hehety/repos/gox/logger/tes/storage.db?cache=shared&_synchronous=0")
	db.DB()
	if err != nil {
		logger.Fatal(err)
	}
	rows, err := db.DB().Query("select id from file")
	if err != nil {
		logger.Error(err)
		return
	}
	cols, _ := rows.Columns()
	logger.Info("columns: ", cols)
	for rows.Next() {
		var id int
		if err = rows.Scan(&id); err != nil {
			logger.Error(err)
			return
		}
		logger.Info("id=", id)
	}
}
