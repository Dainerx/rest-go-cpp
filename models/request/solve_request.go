package request

import (
	"database/sql"
	"time"

	"github.com/Dainerx/rest-go-cpp/pkg/slice"
)

type SolveRequest struct {
	id     int64
	Solver string ""
	Input  string ""
	date   int64
}

var SOLVERS = []string{"loop1e9", "loop1e10", "loop1K", "loop2e9", "loop4e9", "loop10K"}

func NewSolverRequest(solver string, input string) *SolveRequest {
	var sr SolveRequest
	sr.Solver = solver
	sr.Input = input
	sr.date = time.Now().Unix()
	return &sr
}

func (sr SolveRequest) Id() int64 {
	return sr.id
}

func AllSolveRequests(db *sql.DB) ([]*SolveRequest, error) {
	rows, err := db.Query("SELECT * FROM solve_request")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	srs := make([]*SolveRequest, 0)
	for rows.Next() {
		sr := new(SolveRequest)
		err := rows.Scan(&sr.id, &sr.Solver, &sr.Input, &sr.date)
		if err != nil {
			return nil, err
		}
		srs = append(srs, sr)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return srs, nil
}

func GetSolveRequest(db *sql.DB, id int64) (SolveRequest, error) {
	rows, err := db.Query("SELECT * FROM solve_request where id=?", id)
	sr := new(SolveRequest)
	if err != nil {
		return *sr, err
	} else {
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&sr.id, &sr.Solver, &sr.Input, &sr.date)
			if err != nil {
				return *sr, err
			}
		}
		return *sr, nil
	}
}

func AddSolveRequest(db *sql.DB, sr *SolveRequest) error {
	result, err := db.Exec("INSERT INTO solve_request (solver,input,date,user) VALUES(?,?,?,?)", (*sr).Solver, (*sr).Input, (*sr).date)
	if err != nil {
		return err
	}
	sr.id, _ = result.LastInsertId()
	return nil
}

//solver response should be passed here as reference
func (sr SolveRequest) Correct() bool {
	return !sr.empty() && sr.inputCorrect() && sr.solverExists()
}

func (sr SolveRequest) solverExists() bool {
	return slice.ContainString(SOLVERS, sr.Solver)
}

func (sr SolveRequest) empty() bool {
	if sr.Solver == "" || sr.Input == "" {
		return true
	}
	return false
}

func (sr SolveRequest) inputCorrect() bool {
	return true
}
