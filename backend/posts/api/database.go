package api

import (
	"database/sql"
	"log"

	//MySQL driver
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() *sql.DB {

	log.Println("attempting connections")

	var err error
<<<<<<< HEAD
<<<<<<< HEAD
	DB, err = sql.Open("mysql", "root:root@tcp(172.28.1.3:3307)/auth")
=======
=======
>>>>>>> 0f792e0fafe93ecd734de4d058d76046a9c4b1e6
	DB, err = sql.Open("mysql", "root:root@tcp(172.28.1.2:3306)/postsDB?parseTime=true")
	// DB, err = sql.Open("mysql", "root:root@tcp(localhost:3306)/postsDB")



	if err != nil {
		log.Println("couldnt connect")
		panic(err.Error())
	}
<<<<<<< HEAD
>>>>>>> master
=======
>>>>>>> 0f792e0fafe93ecd734de4d058d76046a9c4b1e6

	err = DB.Ping()
	if err != nil {
		log.Println("couldnt ping")
		panic(err.Error())
	}

	return DB
}
