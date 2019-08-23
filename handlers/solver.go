package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Dainerx/rest-go-cpp/cpp"
	"github.com/Dainerx/rest-go-cpp/models"
	"github.com/Dainerx/rest-go-cpp/models/request"
	"github.com/Dainerx/rest-go-cpp/models/response"
)

var UsersRunningInstancesChannels = make(map[int64](chan string))
var UserRunningInstancesRequests = make(map[int64](request.SolveRequest))

const (
	TIME_TO_RUN  = 3 // seconds
	ERROR_SYSTEM = "Error System"
	ERROR_SOLVER = "Error Solver"
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
		return
	}

	// Check wether user has an instance already running.
	timer := time.NewTimer(time.Microsecond * TIME_TO_RUN)
	if _, ok := UsersRunningInstancesChannels[(*ptruser).Id]; ok {
		// If user has a running instance check wether the solver finished or not.
		select {
		case <-UsersRunningInstancesChannels[(*ptruser).Id]:
			// Solver finished return a redirect response telling the user to check
			// /recent or /status to retrive output.
			res := response.RedirectResponse(*ptrreq, response.MESSAGE_INSTANCE_RUNNING_GO, *ptruser)
			json.NewEncoder(w).Encode(res)
		case <-timer.C:
			// Solver did not finish return a redirect response telling the user to check
			// /status for the yet running instance.
			res := response.RedirectResponse(*ptrreq, response.MESSAGE_INSTANCE_RUNNING_WAIT, *ptruser)
			json.NewEncoder(w).Encode(res)
		}
		return
	}

	// Create a timer to wait TIME_TO_RUN seconds before returning wait and redirecting to /status.
	timer = time.NewTimer(time.Second * TIME_TO_RUN)
	// Create a channel and maps it to the current user.
	// Maps the current request to the user in order to persist it.
	crunning_instance := make(chan string)
	UsersRunningInstancesChannels[(*ptruser).Id], UserRunningInstancesRequests[(*ptruser).Id] = crunning_instance, *ptrreq
	go func() { // Run the solver on different subroutine
		output, err := cpp.Run(req.Solver)
		if err != nil {
			UsersRunningInstancesChannels[(*ptruser).Id] <- ERROR_SOLVER // Unlocks channel with solver error.
		} else {
			UsersRunningInstancesChannels[(*ptruser).Id] <- output // Unlocks channel with message.
			err := response.AddSolveResponse(models.Db, response.SuccessResponse(*ptrreq, output, *ptruser))
			if err != nil {
				UsersRunningInstancesChannels[(*ptruser).Id] <- ERROR_SYSTEM // Unlocks channel with system error
			} else {
				fmt.Println("Added a solve response to database")
			}
		}
	}()

	select {
	// Case solver has finished before time out
	case output := <-UsersRunningInstancesChannels[(*ptruser).Id]:
		defer delete(UsersRunningInstancesChannels, (*ptruser).Id)
		defer delete(UserRunningInstancesRequests, (*ptruser).Id)
		if output == ERROR_SYSTEM {
			// System has failed to persist the response in the database
			w.WriteHeader(http.StatusInternalServerError)
			respondInternalServerError(w, *ptrreq, *ptruser)
		} else if output == ERROR_SOLVER {
			// Solver failed and System persisted the response in the database
			res := response.ErrorResponse(*ptrreq, response.MESSAGE_SOLVER_FAILED, *ptruser)
			json.NewEncoder(w).Encode(res)
		} else {
			// Solver finished without error responds with Success
			res := response.SuccessResponse(*ptrreq, output, *ptruser)
			json.NewEncoder(w).Encode(res)
		}
		return
	// Case time out responds wait
	case <-timer.C:
		res := response.WaitingResponse(*ptrreq, *ptruser)
		json.NewEncoder(w).Encode(res)
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

	if crunning_instance, ok := UsersRunningInstancesChannels[(*ptruser).Id]; ok {
		timer := time.NewTimer(time.Second * TIME_TO_RUN)
		select {
		case output := <-crunning_instance:
			defer delete(UsersRunningInstancesChannels, (*ptruser).Id)
			defer delete(UserRunningInstancesRequests, (*ptruser).Id)
			if output == ERROR_SYSTEM {
				w.WriteHeader(http.StatusInternalServerError)
				respondInternalServerError(w, UserRunningInstancesRequests[(*ptruser).Id], *ptruser)
			} else if output == ERROR_SOLVER {
				res := response.ErrorResponse(UserRunningInstancesRequests[(*ptruser).Id], response.MESSAGE_SOLVER_FAILED, *ptruser)
				json.NewEncoder(w).Encode(res)
			} else {
				res := response.SuccessResponse(UserRunningInstancesRequests[(*ptruser).Id], output, *ptruser)
				json.NewEncoder(w).Encode(res)
			}
			return
		case <-timer.C:
			res := response.WaitingResponse(UserRunningInstancesRequests[(*ptruser).Id], *ptruser)
			json.NewEncoder(w).Encode(res)
			return
		}
	}
	res := response.SuccessResponse(UserRunningInstancesRequests[(*ptruser).Id], response.NO_OUTPUT, *ptruser)
	json.NewEncoder(w).Encode(res)
}

func respondInternalServerError(w http.ResponseWriter, sr request.SolveRequest, user models.User) {
	w.WriteHeader(http.StatusInternalServerError)
	res := response.ErrorResponse(sr, response.MESSAGE_INTERNAL_ERROR, user)
	json.NewEncoder(w).Encode(res)
}
