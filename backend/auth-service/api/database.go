package api

import (
	"database/sql"

	"github.com/joho/godotenv"

	//MySQL driver
	_ "github.com/go-sql-driver/mysql"
)

//DB represents the connection to the MySQL database
var (
	DB *sql.DB
)

//InitDB creates the MySQL database connection
func InitDB() *sql.DB {

	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}

	DB, err = sql.Open("mysql", "root:root@/auth")
	if err != nil {
		panic(err.Error())
	}

	err = DB.Ping()
	if err != nil {
		panic(err.Error())
	}

	return DB
}
