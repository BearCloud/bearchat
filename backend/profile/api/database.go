package api

import (
	"database/sql"
	"log"

	"github.com/joho/godotenv"

	//MySQL driver
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
	}

	db, err := sql.Open("mysql", "username:password@tcp([ip]:[port])/profiles")

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()
}
