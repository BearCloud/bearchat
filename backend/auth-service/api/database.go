package api

import (
	"database/sql"
	"log"

	//MySQL driver
	_ "github.com/go-sql-driver/mysql"
)

//DB represents the connection to the MySQL database
var (
	DB *sql.DB
)

//InitDB creates the MySQL database connection
func InitDB() *sql.DB {

	log.Println("attempting connections")

	var err error
	DB, err = sql.Open("mysql", "root:root@tcp(172.28.1.2:3306)/auth")

	if err != nil {
		log.Println("couldnt connect")
		panic(err.Error())
	}

	return DB
}
