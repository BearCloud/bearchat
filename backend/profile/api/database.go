package api

import (
	"database/sql"
	 _ "log"

	//MySQL driver
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() *sql.DB {
	var err error
	DB, err = sql.Open("mysql", "root:root@/profiles")

	if err != nil {
		panic(err.Error())
	}

	return DB
}
