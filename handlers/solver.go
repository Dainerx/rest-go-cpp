package handlers

import (
	"betell-rest/cpp"
	"betell-rest/models/request"
	"betell-rest/models/response"
	"encoding/json"
	"net/http"
	"time"
)

var UsersRunningInstances = make(map[int64](chan string))
var UserRequests = make(map[int64](*request.SolveRequest))

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
		json.NewEncoder(w).Encode(response.UnauthorizedResponse())
		return
	}

	var reqg request.SolveRequest
	json.NewDecoder(r.Body).Decode(&reqg)
	req := request.NewSolverRequest(reqg.Solver, reqg.Input, *user)
	if (*req).Correct() == false {
		res, err := response.ErrorResponse(req, MESSAGE_CHECK_INPUT)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(res)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(res)
		}
		return
	}

	timer := time.NewTimer(time.Second * 1)
	c := make(chan string)
	//test if it exisits already
	UsersRunningInstances[(*user).Id] = c //important
	UserRequests[(*user).Id] = req
	go func() {
		output, err := cpp.Run(req.Solver)
		if err != nil {
			c <- "error" // unblock channel with error
		} else {
			c <- output // unlock channel with message
		}
	}()
	select {
	case output := <-c:
		if output != "error" {
			res, err := response.SuccessResponse(req, output)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			json.NewEncoder(w).Encode(res)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			res, err := response.ErrorResponse(req, MESSAGE_INTERNAL_ERROR)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			json.NewEncoder(w).Encode(res)
		}
		return
	case <-timer.C:
		res, err := response.WaitingResponse(req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(res)
		return
	}
}

func Status(w http.ResponseWriter, r *http.Request) {
	//Declare return content type for the route
	w.Header().Set("Content-Type", "application/json")
	isAuthenticated, user := Authenticated(w, r)
	if isAuthenticated == false {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response.UnauthorizedResponse())
		return
	}

	if c, ok := UsersRunningInstances[(*user).Id]; ok {
		select {
		case <-c:
			delete(UsersRunningInstances, (*user).Id)
			delete(UserRequests, (*user).Id)
			res, err := response.SuccessResponse(UserRequests[(*user).Id], <-c)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			json.NewEncoder(w).Encode(res)
			return
		default:
			res, err := response.WaitingResponse(UserRequests[(*user).Id])
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			json.NewEncoder(w).Encode(res)
			return
		}
	} else {
		res, err := response.SuccessResponse(UserRequests[(*user).Id], "no output")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(res)
	}
}
