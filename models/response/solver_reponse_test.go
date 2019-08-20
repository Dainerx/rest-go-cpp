package response

import (
	"database/sql"
	"os"
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
