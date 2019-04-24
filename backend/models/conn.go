package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var DBSqlite *sqlx.DB
var DBPG *gorm.DB

func ConnectSqlite(driver, url string) *sqlx.DB {
	db, err := sqlx.Open(driver, url)

	if err != nil {
		panic(fmt.Errorf("connect database error: %v", err))
	}

	DBSqlite = db

	return db
}

// Use gOrm
func ConnectPG(url string) *gorm.DB {
	db, err := gorm.Open("postgres", url)

	if err != nil {
		panic(err)
	}

	DBPG = db
	return db
}
