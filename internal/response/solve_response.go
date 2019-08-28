package response

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Dainerx/rest-go-cpp/internal"

	"github.com/Dainerx/rest-go-cpp/internal/request"
)

type SolveResponse struct {
	id            int
	SolverRequest request.Request `json:",omitempty"`
	Status        string          ""
	Message       string          ""
	Output        string          `json:",omitempty"`
	date          int64
	user          internal.User
}

const (
	OK                            = "OK"
	NOT_OK                        = "NOT OK"
	ERROR                         = "ERROR"
	MESSAGE_UNAUTHORIZED          = "Unauthorized, please authenificate on /authenticate"
	MESSAGE_INSTANCE_RUNNING_WAIT = "We are still processing a previous problem, please check your instance's status on /status"
	MESSAGE_INSTANCE_RUNNING_GO   = "We have finished processing your previous problem, please go see your instance's output on /status or /recent"
	MESSAGE_INSTANCE_FINISHED     = "We have finished processing your problem."
	MESSAGE_INSTANCE_NOT_FINISHED = "We are processing your instance, please hold on."
	MESSAGE_NO_INSTANCE           = "You have no solver instance running."
	MESSAGE_CHECK_INPUT           = "Request body is wrong."
	MESSAGE_SOLVER_NOT_FOUND      = "Solver not found!"
	MESSAGE_SOLVER_FAILED         = "Solver failed, no output generated."
	MESSAGE_INTERNAL_ERROR        = "Something went wrong..."
	MESSAGE_LIMIT_10              = "You can only retrive the 10 most recent requests."
	NO_OUTPUT                     = "no output"
	NO_OUTPUT_YET                 = "no output yet"
)

// Return an instance of SolveResponse with Status OK and an output.
func SuccessResponse(sr request.Request, output string, user internal.User) SolveResponse {
	sres := new(SolveResponse)
	(*sres).SolverRequest = sr
	(*sres).Status = OK
	if output != NO_OUTPUT && output != "" {
		(*sres).Message = MESSAGE_INSTANCE_FINISHED
		(*sres).Output = output
	} else {
		(*sres).Message = MESSAGE_NO_INSTANCE
	}
	(*sres).date = time.Now().Unix()
	(*sres).user = user
	return *sres
}

// Return an instance of SolveResponse with Status NOT OK and a message passed as argument.
func RedirectResponse(sr request.Request, message string, user internal.User) SolveResponse {
	sres := new(SolveResponse)
	(*sres).SolverRequest = sr
	(*sres).Status = NOT_OK
	(*sres).Message = message
	(*sres).date = time.Now().Unix()
	(*sres).user = user
	return *sres
}

// Return an instance of SolveResponse with Status OK and MESSAGE_INSTANCE_NOT_FINISHED as message.
func WaitingResponse(sr request.Request, user internal.User) SolveResponse {
	sres := new(SolveResponse)
	(*sres).SolverRequest = sr
	(*sres).Status = OK
	(*sres).Message = MESSAGE_INSTANCE_NOT_FINISHED
	(*sres).date = time.Now().Unix()
	(*sres).user = user
	return *sres
}

// Return an instance of SolveResponse with Status ERROR and a message passed as argument.
func ErrorResponse(sr request.Request, message string, user internal.User) SolveResponse {
	sres := new(SolveResponse)
	(*sres).SolverRequest = sr
	(*sres).Status = ERROR
	(*sres).Message = message
	(*sres).date = time.Now().Unix()
	(*sres).user = user
	return *sres
}

// Return an instance of SolveResponse with Status ERROR and MESSAGE_UNAUTHORIZED as message.
func UnauthorizedResponse() SolveResponse {
	srs := new(SolveResponse)
	(*srs).Status = ERROR
	(*srs).Message = MESSAGE_UNAUTHORIZED
	return *srs
}

// Persist an instance of SolveResponse in the database, and calls request.AddSolveRequest(db *sql.DB, r *Request)
// to persist the request attached to the instance.
// Return nil if the SolveResponse is persisted successfully
func AddSolveResponse(db *sql.DB, sres SolveResponse) error {
	if err := request.AddSolveRequest(db, &sres.SolverRequest); err != nil {
		return err
	}
	_, err := db.Exec("INSERT INTO solve_response (solver_request,status,message,output,date, user) VALUES(?,?,?,?,?,?)", sres.SolverRequest.Id(), sres.Status, sres.Message, sres.Output, sres.date, sres.user.Id)
	if err != nil {
		return err
	}
	return nil
}

func AllSolveResponses(db *sql.DB) ([]SolveResponse, error) {
	var sresponses []SolveResponse
	rows, err := db.Query("SELECT solver_request, status, message, output, date, user FROM solve_response")
	if err != nil {
		return sresponses, err
	}
	defer rows.Close()
	for rows.Next() {
		var idsr int64
		var iduser int64
		var sres SolveResponse
		err := rows.Scan(&idsr, &sres.Status, &sres.Message, &sres.Output, &sres.date, &iduser)
		if err != nil {
			return sresponses, err
		}
		sr, err := request.GetSolveRequest(db, idsr)
		if err != nil {
			return sresponses, err
		}
		sres.SolverRequest = sr
		user, err := internal.GetUser(db, iduser)
		if err != nil {
			return sresponses, err
		}
		sres.user = user

		sresponses = append(sresponses, sres)
	}
	return sresponses, nil
}

func AllSolveResponsesPerUser(db *sql.DB, user internal.User) ([]SolveResponse, error) {
	var sresponsesuser []SolveResponse
	rows, err := db.Query("SELECT solver_request, status, message, output, date FROM solve_response WHERE user=?", user.Id)
	if err != nil {
		fmt.Println(err.Error())
		return sresponsesuser, err
	}
	defer rows.Close()
	for rows.Next() {
		var idsr int64
		var sres SolveResponse
		err := rows.Scan(&idsr, &sres.Status, &sres.Message, &sres.Output, &sres.date)
		if err != nil {
			fmt.Println(err.Error())
			return sresponsesuser, err
		}
		sr, err := request.GetSolveRequest(db, idsr)
		if err != nil {
			fmt.Println(err.Error())
			return sresponsesuser, err
		}
		sres.SolverRequest = sr
		sresponsesuser = append(sresponsesuser, sres)
	}
	return sresponsesuser, nil

}

func RecentSolveResponsesPerUser(db *sql.DB, user internal.User, fetch int8) ([]SolveResponse, error) {
	var sresponses []SolveResponse
	rows, err := db.Query("SELECT solver_request, status, message, output, date FROM solve_response WHERE user=? ORDER BY date DESC LIMIT ?", user.Id, fetch)
	if err != nil {
		return sresponses, err
	}
	defer rows.Close()
	for rows.Next() {
		var idsr int64
		var sres SolveResponse
		err := rows.Scan(&idsr, &sres.Status, &sres.Message, &sres.Output, &sres.date)
		if err != nil {
			return sresponses, err
		}
		sr, err := request.GetSolveRequest(db, idsr)
		if err != nil {
			return sresponses, err
		}
		sres.SolverRequest = sr
		sresponses = append(sresponses, sres)
	}
	return sresponses, nil
}
