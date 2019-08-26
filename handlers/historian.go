package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Dainerx/rest-go-cpp/internal"
	"github.com/Dainerx/rest-go-cpp/internal/response"
	log "github.com/sirupsen/logrus"
)

func Last(w http.ResponseWriter, r *http.Request) {
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
		"func": "historian.Last",
	})
	sresponses, err := response.RecentSolveResponses(internal.Db, *ptruser, 1)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Errorf("RecentSolveResponses failed: %v.", err)
		return
	}
	json.NewEncoder(w).Encode(sresponses[0])
	logger.Tracef("RecentSolveResponses succeeded with limit 1.")
}

func Recent(w http.ResponseWriter, r *http.Request) {
	// Returning json response
	w.Header().Set("Content-Type", "application/json")
	//Checks if user is authenificated.
	isAuthenticated, ptruser := Authenticated(w, r)
	if isAuthenticated == false {
		// Responds Unauthorized.
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// init logger
	logger := log.WithFields(log.Fields{
		"user": (*ptruser).Email,
		"func": "historian.Recent",
	})

	// Parse fetch variable from /recent?fetch=x
	queryparm := r.URL.Query().Get("fetch")
	if queryparm == "" {
		// Param fetch not found in the url
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.ErrorResponse(nil, "Fetch param not found.", *ptruser))
		logger.Tracef("Failed to retrive query param fetch.")
		return
	}
	fetch, err := strconv.ParseInt(queryparm, 10, 8)
	if err != nil {
		// Error parsing string to int
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.ErrorResponse(nil, "Fetch should be of type integer.", *ptruser))
		logger.Tracef("Failed to parse fetch: %v", err)
		return
	}
	if fetch > 10 {
		// Breach in fetching
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode(response.ErrorResponse(nil, response.MESSAGE_LIMIT_10, *ptruser))
		logger.Tracef("Failed to fetch, want 10 or less have %d.", fetch)
		return
	}

	sresponses, err := response.RecentSolveResponses(internal.Db, *ptruser, int8(fetch))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Errorf("RecentSolveResponses failed: %v.", err)
		return
	}
	json.NewEncoder(w).Encode(sresponses)
	logger.Tracef("RecentSolveResponses succeeded with fetch %d.", fetch)
}
