package api

import (
	"database/sql"
	_"log"

	//MySQL driver
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() *sql.DB {

	var err error
	DB, err = sql.Open("mysql", "root:root@tcp(192.168.50.166:3306)/postsDB?parseTime=true")

	if err != nil {
		panic(err.Error())
	}

	return DB
}
