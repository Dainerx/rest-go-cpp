package models

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB

func InitDB() {
	user := os.Getenv("DB_BETELL_USER")
	pass := os.Getenv("DB_BETELL_PASS")
	database := os.Getenv("DB_BETELL_DB")
	var err error
	Db, err = sql.Open("mysql", user+":"+pass+"@/"+database)
	if err != nil {
		log.Panic(err)
	}

	if err = Db.Ping(); err != nil {
		log.Panic(err)
	}
}
