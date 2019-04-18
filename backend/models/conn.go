package models

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sqlx.DB

func ConnectDatabase(driver, url string) *sqlx.DB {
	db, err := sqlx.Open(driver, url)

	if err != nil {
		panic(fmt.Errorf("connect database error: %v", err))
	}

	DB = db

	return db
}
