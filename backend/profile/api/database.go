package api

import (
	"database/sql"
	"log"

	"github.com/joho/godotenv"

	//MySQL driver
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
	}

	DB, err = sql.Open("mysql", "root:root@/profiles")

	if err != nil {
		panic(err.Error())
	}

	defer DB.Close()
}
