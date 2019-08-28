package internal

import (
	"database/sql"
	"os"
	"strconv"
	"testing"

	txdb "github.com/DATA-DOG/go-txdb"

	"github.com/romanyx/polluter"
)

func init() {
	user := os.Getenv("DB_BETELL_USER")
	pass := os.Getenv("DB_BETELL_PASS")
	database := os.Getenv("DB_BETELL_DB")
	txdb.Register("txdb", "mysql", user+":"+pass+"@/"+database)
}

func openDb() (*sql.DB, error, func() error) {
	user := os.Getenv("DB_BETELL_USER")
	pass := os.Getenv("DB_BETELL_PASS")
	database := os.Getenv("DB_BETELL_DB")
	var err error
	db, err := sql.Open("txdb", user+":"+pass+"@/"+database)
	return db, err, db.Close
}

func polluteDb(db *sql.DB, t *testing.T) {
	seed, err := os.Open(os.Getenv("DB_BETELL_SEED"))
	if err != nil {
		t.Fatalf("failed to open seed file: %s", err)
	}
	defer seed.Close()
	p := polluter.New(polluter.MySQLEngine(db))
	if err := p.Pollute(seed); err != nil {
		t.Fatalf("failed to pollute: %s", err)
	}
}

func TestAllUsers(t *testing.T) {
	db, err, closer := openDb()
	if err != nil {
		t.Fatal(err)
	}
	defer closer()

	polluteDb(db, t)

	users, err := AllUsers(db)
	if err != nil {
		t.Fatalf("AllUsers() failed: %v", err)
	}
	got := users[0].Email
	if got != "dainer@gmail.com" {
		t.Errorf("AllUsers().Email = %s; want dainer@gmail.com", got)
	}
}

func TestAddUser(t *testing.T) {
	db, err, closer := openDb()
	if err != nil {
		t.Fatal(err)
	}
	defer closer()

	users, _ := AllUsers(db)
	userscount := len(users)
	err = AddUser(db, "bla@gmail.com", "pass")
	if err != nil {
		t.Fatalf("AddUser() failed: %v", err)
	}

	users, _ = AllUsers(db)
	if (userscount + 1) != len(users) {
		t.Fatal("AddUser() failed to write in database")
	}
}

func TestUserExists(t *testing.T) {
	db, err, closer := openDb()
	if err != nil {
		t.Fatal(err)
	}
	defer closer()

	_ = AddUser(db, "test@gmail.com", "pass")
	got, _, err := UserExists(db, "test@gmail.com", "pass")
	if err != nil {
		t.Errorf("UserExists(test@gmail.com,pass) failed: %v", err)
	}
	if got != true {
		t.Errorf("UserExists(test@gmail.com,pass) = %s; want true", strconv.FormatBool(got))
	}

	got, _, err = UserExists(db, "1@gmail.com", "pass")
	if err != nil {
		t.Errorf("UserExists(1@gmail.com,pass) failed")
	}
	if got != false {
		t.Errorf("UserExists(1@gmail.com,pass) = %s; want false", strconv.FormatBool(got))
	}

}
