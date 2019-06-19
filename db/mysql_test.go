package db_test

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hetianyi/gox/logger"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"testing"
)

func init() {
	logger.Init(nil)
}

type Article struct {
	Title string `gorm:"column:ma_title"`
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

func TestPool(t *testing.T) {

}
