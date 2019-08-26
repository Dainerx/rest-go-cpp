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

func RedirectResponse(sr request.Request, message string, user internal.User) SolveResponse {
	sres := new(SolveResponse)
	(*sres).SolverRequest = sr
	(*sres).Status = OK
	(*sres).Message = message
	(*sres).date = time.Now().Unix()
	(*sres).user = user
	return *sres
}

func WaitingResponse(sr request.Request, user internal.User) SolveResponse {
	sres := new(SolveResponse)
	(*sres).SolverRequest = sr
	(*sres).Status = OK
	(*sres).Message = MESSAGE_INSTANCE_NOT_FINISHED
	(*sres).date = time.Now().Unix()
	(*sres).user = user
	return *sres
}

func ErrorResponse(sr request.Request, message string, user internal.User) SolveResponse {
	sres := new(SolveResponse)
	(*sres).SolverRequest = sr
	(*sres).Status = ERROR
	(*sres).Message = message
	(*sres).date = time.Now().Unix()
	(*sres).user = user
	return *sres
}

func UnauthorizedResponse() SolveResponse {
	srs := new(SolveResponse)
	(*srs).Status = ERROR
	(*srs).Message = MESSAGE_UNAUTHORIZED
	return *srs
}

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

func RecentSolveResponses(db *sql.DB, user internal.User, fetch int8) ([]SolveResponse, error) {
	var sresponses []SolveResponse
	rows, err := db.Query("SELECT solver_request, status, message, output, date FROM solve_response WHERE user=? ORDER BY date DESC LIMIT ?", user.Id, fetch)
	if err != nil {
		fmt.Println(err.Error())
		return sresponses, err
	}
	defer rows.Close()
	for rows.Next() {
		var idsr int64
		var sres SolveResponse
		err := rows.Scan(&idsr, &sres.Status, &sres.Message, &sres.Output, &sres.date)
		if err != nil {
			fmt.Println(err.Error())
			return sresponses, err
		}
		sr, err := request.GetSolveRequest(db, idsr)
		if err != nil {
			fmt.Println(err.Error())
			return sresponses, err
		}
		sres.SolverRequest = sr
		sresponses = append(sresponses, sres)
	}
	return sresponses, nil
}
