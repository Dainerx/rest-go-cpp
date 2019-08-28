package internal

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB

func InitDB() error {
	user := os.Getenv("DB_BETELL_USER")
	pass := os.Getenv("DB_BETELL_PASS")
	database := os.Getenv("DB_BETELL_DB")
	var err error
	Db, err = sql.Open("mysql", user+":"+pass+"@/"+database)
	if err != nil {
		return err
	}
	if err = Db.Ping(); err != nil {
		return err
	}
	return nil
}
