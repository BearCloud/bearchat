package api

import (
	"database/sql"
	"github.com/joho/godotenv"
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

	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}

	log.Println("attempting connections")

	DB, err = sql.Open("mysql", "root:root@tcp(172.28.1.2:3306)/auth")

	if err != nil {
		log.Println("couldnt connect")
		panic(err.Error())
	}

	err = DB.Ping()
	if err != nil {
		log.Println("couldnt ping")
		panic(err.Error())
	}

	return DB
}
