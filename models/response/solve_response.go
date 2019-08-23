package response

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Dainerx/rest-go-cpp/models"
	"github.com/Dainerx/rest-go-cpp/models/request"
)

type SolveResponse struct {
	id            int
	solverRequest request.SolveRequest
	Status        string ""
	Message       string ""
	Output        string
	date          int64
	user          models.User
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
	NO_OUTPUT                     = "no output"
	NO_OUTPUT_YET                 = "no output yet"
)

func SuccessResponse(sr request.SolveRequest, output string, user models.User) SolveResponse {
	sres := new(SolveResponse)
	(*sres).solverRequest = sr
	(*sres).Status = OK
	if output != NO_OUTPUT && output != "" {
		(*sres).Message = MESSAGE_INSTANCE_FINISHED
		(*sres).Output = output
	} else {
		(*sres).Message = MESSAGE_NO_INSTANCE
		(*sres).Output = NO_OUTPUT
	}
	(*sres).date = time.Now().Unix()
	(*sres).user = user
	return *sres
}

func RedirectResponse(sr request.SolveRequest, message string, user models.User) SolveResponse {
	sres := new(SolveResponse)
	(*sres).solverRequest = sr
	(*sres).Status = OK
	(*sres).Message = message
	(*sres).Output = NO_OUTPUT
	(*sres).date = time.Now().Unix()
	(*sres).user = user
	return *sres
}

func WaitingResponse(sr request.SolveRequest, user models.User) SolveResponse {
	sres := new(SolveResponse)
	(*sres).solverRequest = sr
	(*sres).Status = OK
	(*sres).Message = MESSAGE_INSTANCE_NOT_FINISHED
	(*sres).Output = NO_OUTPUT_YET
	(*sres).date = time.Now().Unix()
	(*sres).user = user
	return *sres
}

func ErrorResponse(sr request.SolveRequest, message string, user models.User) SolveResponse {
	sres := new(SolveResponse)
	(*sres).solverRequest = sr
	(*sres).Status = ERROR
	(*sres).Message = message
	(*sres).Output = NO_OUTPUT
	(*sres).date = time.Now().Unix()
	(*sres).user = user
	return *sres
}

func UnauthorizedResponse() SolveResponse {
	srs := new(SolveResponse)
	(*srs).Status = ERROR
	(*srs).Message = MESSAGE_UNAUTHORIZED
	(*srs).Output = NO_OUTPUT
	return *srs
}

func AddSolveResponse(db *sql.DB, sres SolveResponse) error {
	if err := request.AddSolveRequest(db, &sres.solverRequest); err != nil {
		return err
	}
	_, err := db.Exec("INSERT INTO solve_response (solver_request,status,message,output,date) VALUES(?,?,?,?,?)", sres.solverRequest.Id(), sres.Status, sres.Message, sres.Output, sres.date)
	if err != nil {
		return err
	}
	return nil
}

func RecentSolveResponses(db *sql.DB, user models.User, limit int) ([]SolveResponse, error) {
	var sresponses []SolveResponse
	rows, err := db.Query("SELECT solver_request, status, message, output, date FROM solve_response WHERE user=? ORDER BY date DESC LIMIT ?", user.Id, limit)
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
		sres.solverRequest = sr
		sresponses = append(sresponses, sres)
	}
	return sresponses, nil
}
