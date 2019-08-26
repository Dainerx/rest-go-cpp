package request

import (
	"database/sql"
	"os"
	"strconv"
	"testing"

	txdb "github.com/DATA-DOG/go-txdb"
	_ "github.com/go-sql-driver/mysql"
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

func TestSolveRequestTypeRequest(t *testing.T) {
	var _ Request = SolveRequest{}       // Verify that SolveRequest implements Request.
	var _ Request = (*SolveRequest)(nil) // Verify that *SolveRequest implements Request.
}
func TestAllSolveRequests(t *testing.T) {
	db, err, closer := openDb()
	if err != nil {
		t.Fatal(err)
	}
	defer closer()

	polluteDb(db, t)

	_, err = AllSolveRequests(db)
	if err != nil {
		t.Fatalf("AllSolveRequests() failed:%s", err)
	}
}

func TestAddSolveRequest(t *testing.T) {
	db, err, closer := openDb()
	if err != nil {
		t.Fatal(err)
	}
	defer closer()

	polluteDb(db, t)

	srs, _ := AllSolveRequests(db)
	srscount := len(srs)
	var sr Request = SolveRequest{Solver: "solver", Input: "inputTestSr"} // Verify that SolveRequest implements Request.
	err = AddSolveRequest(db, &sr)
	if err != nil {
		t.Fatalf("AddSolveRequest() failed: %s", err)
	}

	srs, _ = AllSolveRequests(db)
	if (srscount + 1) != len(srs) {
		t.Error("AddSolveRequest() failed to write in database")
	}

	sr1, err := GetSolveRequest(db, sr.Id())
	if err != nil {
		t.Fatalf("GetSolveRequest() failed: %s", err)
	}
	got := sr1.Input
	if got != "inputTestSr" {
		t.Errorf("got %s, want inputTestSr", got)
	}
}

func TestCorrect(t *testing.T) {
	var request SolveRequest
	request.Input = ""
	request.Solver = ""
	got := request.Correct()
	if got != false {
		t.Errorf("request.Correct() = %s; want false", strconv.FormatBool(got))
	}

	request.Input = "my_input"
	request.Solver = "no solver"
	got = request.Correct()
	if got != false {
		t.Errorf("request.Correct() = %s; want false", strconv.FormatBool(got))
	}

	request.Input = ""
	request.Solver = "loop1K"
	got = request.Correct()
	if got != false {
		t.Errorf("request.Correct() = %s; want false", strconv.FormatBool(got))
	}

	request.Input = "my_input"
	request.Solver = "loop1K"
	got = request.Correct()
	if got != true {
		t.Errorf("request.Correct() = %s; want false", strconv.FormatBool(got))
	}
}
