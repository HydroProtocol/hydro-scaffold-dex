package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

var DB *gorm.DB

func Connect(url string) *gorm.DB {
	db, err := gorm.Open("postgres", url)

	if err != nil {
		panic(err)
	}

	gorm.NowFunc = func() time.Time {
		return time.Now().UTC()
	}

	DB = db
	return db
}
