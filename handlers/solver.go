package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Dainerx/rest-go-cpp/cpp"
	"github.com/Dainerx/rest-go-cpp/internal"
	"github.com/Dainerx/rest-go-cpp/internal/request"
	"github.com/Dainerx/rest-go-cpp/internal/response"
	log "github.com/sirupsen/logrus"
)

var UsersRunningInstancesChannels = make(map[int64](chan string))
var UserRunningInstancesRequests = make(map[int64](request.SolveRequest))

const (
	SOLVER_RUN_TIMEOUT    = 5  // seconds
	SOLVER_STATUS_TIMEOUT = 10 //seconds
	ERROR_SOLVER          = "Error Solver: "
)

func Solve(w http.ResponseWriter, r *http.Request) {
	// Returning json response
	w.Header().Set("Content-Type", "application/json")
	//Checks if user is authenificated.
	isAuthenticated, ptruser := Authenticated(w, r)
	if isAuthenticated == false {
		// Responds Unauthorized.
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response.UnauthorizedResponse())
		return
	}
	//init logger
	logger := log.WithFields(log.Fields{
		"user": (*ptruser).Email,
		"func": "solver.Solve",
	})
	// Recieve SolveRequest decoded it and check wether it is correct or not.
	var req request.SolveRequest
	json.NewDecoder(r.Body).Decode(&req)
	solver, input := req.Solver, req.Input
	ptrreq := request.NewSolverRequest(solver, input)
	if (*ptrreq).Correct() == false {
		// Responds Error with input.
		res := response.ErrorResponse(*ptrreq, response.MESSAGE_CHECK_INPUT, *ptruser)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		logger.Trace("Request is not correct.")
		return
	}

	// Check wether user has an instance already running.
	timer := time.NewTimer(time.Microsecond * SOLVER_RUN_TIMEOUT)
	if _, ok := UsersRunningInstancesChannels[(*ptruser).Id]; ok {
		// If user has a running instance check wether the solver finished or not.
		select {
		case <-UsersRunningInstancesChannels[(*ptruser).Id]:
			// Solver finished return a redirect response telling the user to check
			// /recent or /status to retrive output.
			res := response.RedirectResponse(*ptrreq, response.MESSAGE_INSTANCE_RUNNING_GO, *ptruser)
			json.NewEncoder(w).Encode(res)
			logger.Trace("Found instance finished running.")
		case <-timer.C:
			// Solver did not finish return a redirect response telling the user to check
			// /status for the yet running instance.
			res := response.RedirectResponse(*ptrreq, response.MESSAGE_INSTANCE_RUNNING_WAIT, *ptruser)
			json.NewEncoder(w).Encode(res)
			logger.Trace("Found instance still running.")
		}
		return
	}

	// Create a timer to wait TIME_TO_RUN seconds before returning wait and redirecting to /status.
	timer = time.NewTimer(time.Second * SOLVER_RUN_TIMEOUT)
	// Create a channel and maps it to the current user.
	// Maps the current request to the user in order to persist it.
	crunning_instance := make(chan string)
	UsersRunningInstancesChannels[(*ptruser).Id], UserRunningInstancesRequests[(*ptruser).Id] = crunning_instance, *ptrreq
	go func() { // Run the solver on different subroutine
		output, err := cpp.Run(req.Solver)
		if err != nil {
			errmsg := ERROR_SOLVER + err.Error()
			UsersRunningInstancesChannels[(*ptruser).Id] <- errmsg // Unlocks channel with solver error.
		} else {
			UsersRunningInstancesChannels[(*ptruser).Id] <- output // Unlocks channel with message.
		}
	}()

	select {
	// Case solver has finished before time out
	case output := <-UsersRunningInstancesChannels[(*ptruser).Id]:
		defer delete(UsersRunningInstancesChannels, (*ptruser).Id)
		defer delete(UserRunningInstancesRequests, (*ptruser).Id)
		if strings.Contains(output, ERROR_SOLVER) {
			// Solver has failed.
			w.WriteHeader(http.StatusInternalServerError)
			respondInternalServerError(w, *ptrreq, *ptruser)
			logger.Error(output)
		} else {
			// Solver finished without error responds with Success
			res := response.SuccessResponse(*ptrreq, output, *ptruser)
			err := response.AddSolveResponse(internal.Db, res)
			if err != nil {
				json.NewEncoder(w).Encode(res)
				logger.Errorf("Failed to persist SolveResponse: %v", err)
			}
			json.NewEncoder(w).Encode(res)
			logger.Infof("Finished running the solver %s against input: %s", (*ptrreq).Solver, (*ptrreq).Input)
		}
		return
	// Case time out responds wait
	case <-timer.C:
		res := response.WaitingResponse(*ptrreq, *ptruser)
		json.NewEncoder(w).Encode(res)
		logger.Tracef("Time out to solve for solver %s against input: %s", (*ptrreq).Solver, (*ptrreq).Input)
	}
}

func Status(w http.ResponseWriter, r *http.Request) {
	// Returning json response
	w.Header().Set("Content-Type", "application/json")
	//Checks if user is authenificated.
	isAuthenticated, ptruser := Authenticated(w, r)
	if isAuthenticated == false {
		// Responds Unauthorized.
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response.UnauthorizedResponse())
		return
	}

	//init logger
	logger := log.WithFields(log.Fields{
		"user": (*ptruser).Email,
		"func": "solver.Solve",
	})
	if crunning_instance, ok := UsersRunningInstancesChannels[(*ptruser).Id]; ok {
		timer := time.NewTimer(time.Second * SOLVER_STATUS_TIMEOUT)
		select {
		case output := <-crunning_instance:
			defer delete(UsersRunningInstancesChannels, (*ptruser).Id)
			defer delete(UserRunningInstancesRequests, (*ptruser).Id)
			if strings.Contains(output, ERROR_SOLVER) {
				// Solver has failed.
				w.WriteHeader(http.StatusInternalServerError)
				respondInternalServerError(w, UserRunningInstancesRequests[(*ptruser).Id], *ptruser)
				logger.Error(output)
			} else {
				// Solver finished without error responds with Success
				res := response.SuccessResponse(UserRunningInstancesRequests[(*ptruser).Id], output, *ptruser)
				err := response.AddSolveResponse(internal.Db, res)
				if err != nil {
					json.NewEncoder(w).Encode(res)
					logger.Errorf("Failed to persist SolveResponse: %v", err)
				}
				json.NewEncoder(w).Encode(res)
				logger.Infof("Finished running the solver %s against input: %s", UserRunningInstancesRequests[(*ptruser).Id].Solver, UserRunningInstancesRequests[(*ptruser).Id].Input)
			}
			return
		case <-timer.C:
			res := response.WaitingResponse(UserRunningInstancesRequests[(*ptruser).Id], *ptruser)
			json.NewEncoder(w).Encode(res)
			logger.Tracef("Time out to solve for solver %s against input: %s", UserRunningInstancesRequests[(*ptruser).Id].Solver, UserRunningInstancesRequests[(*ptruser).Id].Input)
			return
		}
	} else {
		res := response.SuccessResponse(UserRunningInstancesRequests[(*ptruser).Id], response.NO_OUTPUT, *ptruser)
		json.NewEncoder(w).Encode(res)
		logger.Trace("No instance running")
	}
}

func respondInternalServerError(w http.ResponseWriter, sr request.SolveRequest, user internal.User) {
	w.WriteHeader(http.StatusInternalServerError)
	res := response.ErrorResponse(sr, response.MESSAGE_INTERNAL_ERROR, user)
	json.NewEncoder(w).Encode(res)
}
