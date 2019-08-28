package response

import (
	"database/sql"
	"os"
	"strings"
	"testing"

	"github.com/Dainerx/rest-go-cpp/internal/request"

	"github.com/Dainerx/rest-go-cpp/internal"

	txdb "github.com/DATA-DOG/go-txdb"
	"github.com/romanyx/polluter"
)

const OUTPUT_TEST = "test output"

type TB interface {
	Fatalf(format string, args ...interface{})
}

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

func polluteDb(db *sql.DB, tb TB) {
	seed, err := os.Open(os.Getenv("DB_BETELL_SEED"))
	if err != nil {
		tb.Fatalf("failed to open seed file: %s", err)
	}
	defer seed.Close()
	p := polluter.New(polluter.MySQLEngine(db))
	if err := p.Pollute(seed); err != nil {
		tb.Fatalf("failed to pollute: %s", err)
	}
}

func TestSucessResponse(t *testing.T) {
	db, err, closer := openDb()
	if err != nil {
		t.Fatal(err)
	}
	defer closer()

	polluteDb(db, t)

	user, _ := internal.GetUser(db, 1)
	req, _ := request.GetSolveRequest(db, 1)
	res := SuccessResponse(req, OUTPUT_TEST, user)

	got := res.Status
	if strings.Compare(got, OK) != 0 {
		t.Fatalf("res.Status failed got %s, want %s", got, OK)
	}
	got = res.Output
	if strings.Compare(got, OUTPUT_TEST) != 0 {
		t.Fatalf("res.Output failed got %s, want %s", got, OUTPUT_TEST)
	}
	if res.user != user {
		t.Fatalf("res.user failed got %v, want %v", res.user, user)
	}
}

func TestErrorResponse(t *testing.T) {
	db, err, closer := openDb()
	if err != nil {
		t.Fatal(err)
	}
	defer closer()

	polluteDb(db, t)

	user, _ := internal.GetUser(db, 1)
	req, _ := request.GetSolveRequest(db, 1)
	res := ErrorResponse(req, MESSAGE_SOLVER_FAILED, user)

	got := res.Status
	if strings.Compare(got, ERROR) != 0 {
		t.Fatalf("res.Status failed got %s, want %s", got, ERROR)
	}
	got = res.Message
	if strings.Compare(got, MESSAGE_SOLVER_FAILED) != 0 {
		t.Fatalf("res.Message failed got %s, want %s", got, MESSAGE_SOLVER_FAILED)
	}
	if res.user != user {
		t.Fatalf("res.user failed got %v, want %v", res.user, user)
	}
}

func TestAddSolveResponse(t *testing.T) {
	db, err, closer := openDb()
	if err != nil {
		t.Fatal(err)
	}
	defer closer()

	polluteDb(db, t)

	user, _ := internal.GetUser(db, 1)
	req, _ := request.GetSolveRequest(db, 1)
	res := SuccessResponse(req, OUTPUT_TEST, user)
	err = AddSolveResponse(db, res)
	if err != nil {
		t.Fatalf("AddSolveResponse(db,%v) failed: %s", res, err)
	}
}

func TestAllSolveResponses(t *testing.T) {
	db, err, closer := openDb()
	if err != nil {
		t.Fatal(err)
	}
	defer closer()

	polluteDb(db, t)

	_, err = AllSolveResponses(db)
	if err != nil {
		t.Fatalf("AllSolveResponses(db) failed: %s", err)
	}
}

func TestAllSolveResponsesPerUser(t *testing.T) {
	db, err, closer := openDb()
	if err != nil {
		t.Fatal(err)
	}
	defer closer()

	polluteDb(db, t)

	user, _ := internal.GetUser(db, 1)
	_, err = AllSolveResponsesPerUser(db, user)
	if err != nil {
		t.Fatalf("AllSolveResponsesPerUser(db,%v) failed: %s", user, err)
	}
}

func TestRecentSolveResponsesPerUser(t *testing.T) {
	db, err, closer := openDb()
	if err != nil {
		t.Fatal(err)
	}
	defer closer()

	polluteDb(db, t)

	user, _ := internal.GetUser(db, 1)
	var fetch int8 = 1
	sresponses, err := RecentSolveResponsesPerUser(db, user, fetch)
	if err != nil {
		t.Fatalf("RecentSolveResponsesPerUser(db,%v) failed: %s", user, err)
	}

	got := len(sresponses)
	if got != int(fetch) {
		t.Fatalf("RecentSolveResponsesPerUser(db,%v,%d) returned slice of length = %d ; want %d", user, fetch, got, fetch)
	}

	fetch++
	sresponses, err = RecentSolveResponsesPerUser(db, user, fetch)
	if err != nil {
		t.Fatalf("RecentSolveResponsesPerUser(db,%v,%d) failed: %s", user, fetch, err)
	}

	got = len(sresponses)
	if got != int(fetch) {
		t.Fatalf("RecentSolveResponsesPerUser(db,%v,%d) returned slice of length = %d ; want %d", user, fetch, got, fetch)
	}

}

func BenchmarkTenRecentSolveResponsesPerUser(b *testing.B) {
	db, err, closer := openDb()
	if err != nil {

		b.Fatal(err)
	}
	defer closer()

	polluteDb(db, b)

	user, _ := internal.GetUser(db, 1)
	var fetch int8 = 10

	for i := 0; i < b.N; i++ {
		rs, _ := RecentSolveResponsesPerUser(db, user, fetch)
		var sr SolveResponse
		rs = append(rs, sr)
	}
}

func BenchmarkAllSolveResponsesPerUser(b *testing.B) {
	db, err, closer := openDb()
	if err != nil {

		b.Fatal(err)
	}
	defer closer()

	polluteDb(db, b)

	user, _ := internal.GetUser(db, 1)

	for i := 0; i < b.N; i++ {
		rs, _ := AllSolveResponsesPerUser(db, user)
		var sr SolveResponse
		rs = append(rs, sr)
	}
}
