package models

import (
	"betell-rest/util"
	"time"
)

type SolveRequest struct {
	id     int64
	Solver string ""
	Input  string ""
	date   int64
	user   User
}

var SOLVERS = []string{"loop1e9", "loop1e10", "loop1K", "loop2e9", "loop4e9", "loop10K"}

func NewSolverRequest(solver string, input string, user User) *SolveRequest {
	var sr SolveRequest
	sr.Solver = solver
	sr.Input = input
	sr.date = time.Now().Unix()
	sr.user = user
	return &sr
}

func addSolveRequest(sr *SolveRequest) error {
	result, err := db.Exec("INSERT INTO solve_request (solver,input,date,user) VALUES(?,?,?,?)", (*sr).Solver, (*sr).Input, (*sr).date, (*sr).user.Id)
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
	return util.ContainString(SOLVERS, sr.Solver)
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
