package handlers

import (
	"betell-rest/cpp"
	"betell-rest/models"
	"encoding/json"
	"net/http"
	"time"
)

var UsersRunningInstances = make(map[int64](chan string))
var UserRequests = make(map[int64](*models.SolveRequest))

const (
	MESSAGE_CHECK_INPUT      = "Request body is wrong."
	MESSAGE_SOLVER_NOT_FOUND = "Solver not found!"
	MESSAGE_INTERNAL_ERROR   = "Something went wrong..."
)

func Solve(w http.ResponseWriter, r *http.Request) {
	//Declare return content type for the route
	w.Header().Set("Content-Type", "application/json")
	isAuthenticated, user := Authenticated(w, r)
	if isAuthenticated == false {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(models.UnauthorizedResponse())
		return
	}

	var re models.SolveRequest
	json.NewDecoder(r.Body).Decode(&re)
	request := models.NewSolverRequest(re.Solver, re.Input, *user)
	if (*request).Correct() == false {
		response, err := models.ErrorResponse(request, MESSAGE_CHECK_INPUT)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
		}
		return
	}

	timer := time.NewTimer(time.Second * 1)
	c := make(chan string)
	//test if it exisits already
	UsersRunningInstances[(*user).Id] = c //important
	UserRequests[(*user).Id] = request
	go func() {
		output, err := cpp.Run(request.Solver)
		if err != nil {
			c <- "error" // unblock channel with error
		} else {
			c <- output // unlock channel with message
		}
	}()
	select {
	case output := <-c:
		if output != "error" {
			response, err := models.SuccessResponse(request, output)
			json.NewEncoder(w).Encode()
		} else {
			w.WriteHeader(http.StatusBadRequest)
			response, err := models.ErrorResponse(request, MESSAGE_INTERNAL_ERROR)
			json.NewEncoder(w).Encode()
		}
		return
	case <-timer.C:
		response, err := models.WaitingResponse(request)
		json.NewEncoder(w).Encode()
		return
	}
}

func Status(w http.ResponseWriter, r *http.Request) {
	//Declare return content type for the route
	w.Header().Set("Content-Type", "application/json")
	isAuthenticated, user := Authenticated(w, r)
	if isAuthenticated == false {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(models.UnauthorizedResponse())
		return
	}

	if c, ok := UsersRunningInstances[(*user).Id]; ok {
		select {
		case <-c:
			delete(UsersRunningInstances, (*user).Id)
			delete(UserRequests, (*user).Id)
			response, err := models.SuccessResponse(UserRequests[(*user).Id], <-c)
			json.NewEncoder(w).Encode()
			return
		default:
			response, err := models.WaitingResponse(UserRequests[(*user).Id])
			json.NewEncoder(w).Encode()
			return
		}
	} else {
		response, err := models.SuccessResponse(UserRequests[(*user).Id], "no output")
		json.NewEncoder(w).Encode()
	}
}
