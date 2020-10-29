package api

import (
	"database/sql"
	"log"
	"time"

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
	// DB, err = sql.Open("mysql", "root:root@tcp(localhost:3306)/auth")


	//Execute a test query
	_, err = DB.Query("SELECT * FROM users")
	for err != nil {
		log.Println("couldnt connect, waiting 20 seconds before retrying")
		time.Sleep(20*time.Second)
		DB, err = sql.Open("mysql", "root:root@tcp(172.28.1.2:3306)/auth")
	}

	return DB
}