package response

import (
	"time"

	"github.com/Dainerx/rest-go-cpp/models"
	"github.com/Dainerx/rest-go-cpp/models/request"
)

type SolveResponse struct {
	id            int
	solverRequest *request.SolveRequest
	Status        string ""
	Message       string ""
	Output        string
	date          int64
}

const (
	OK                        = "OK"
	ERROR                     = "ERROR"
	MESSAGE_UNAUTHORIZED      = "Unauthorized, please authenificate on /authenticate"
	MESSAGE_INSTANCE_FINISHED = "We have finished processing your problem."
	MESSAGE_INSTANCE_RUNNING  = "We are processing your instance, please hold on."
	MESSAGE_NO_INSTANCE       = "You have no solver instance running."
	MESSAGE_INTERNAL_ERROR    = "Something went wrong..."
	NO_OUTPUT                 = "no output"
	NO_OUTPUT_YET             = "no output yet"
)

func SuccessResponse(sr *request.SolveRequest, output string) (SolveResponse, error) {
	srs := new(SolveResponse)
	(*srs).solverRequest = sr
	(*srs).Status = OK
	if output != NO_OUTPUT && output != "" {
		(*srs).Message = MESSAGE_INSTANCE_FINISHED
		(*srs).Output = output
	} else {
		(*srs).Message = MESSAGE_NO_INSTANCE
		(*srs).Output = NO_OUTPUT
	}
	(*srs).date = time.Now().Unix()
	if err := addSolveResponse(srs); err != nil {
		return *srs, err
	}
	return *srs, nil
}

func WaitingResponse(sr *request.SolveRequest) (SolveResponse, error) {
	srs := new(SolveResponse)
	(*srs).solverRequest = sr
	(*srs).Status = OK
	(*srs).Message = MESSAGE_INSTANCE_RUNNING
	(*srs).Output = NO_OUTPUT_YET
	(*srs).date = time.Now().Unix()
	if err := addSolveResponse(srs); err != nil {
		return *srs, err
	}
	return *srs, nil
}

func ErrorResponse(sr *request.SolveRequest, message string) (SolveResponse, error) {
	srs := new(SolveResponse)
	(*srs).solverRequest = sr
	(*srs).Status = ERROR
	(*srs).Message = message
	(*srs).Output = NO_OUTPUT
	(*srs).date = time.Now().Unix()
	if err := addSolveResponse(srs); err != nil {
		return *srs, err
	}
	return *srs, nil
}

func UnauthorizedResponse() SolveResponse {
	srs := new(SolveResponse)
	(*srs).Status = ERROR
	(*srs).Message = MESSAGE_UNAUTHORIZED
	(*srs).Output = NO_OUTPUT
	return *srs
}

func addSolveResponse(srs *SolveResponse) error {
	if err := request.AddSolveRequest((*srs).solverRequest); err != nil {
		return err
	}
	_, err := models.Db.Exec("INSERT INTO solve_request (solver_request,status,message,output,date) VALUES(?,?,?,?,?)", (*srs).solverRequest.GetId(), (*srs).Status, (*srs).Message, (*srs).Output, (*srs).date)
	if err != nil {
		return err
	}
	return nil
}
