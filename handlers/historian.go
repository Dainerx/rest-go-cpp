package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Dainerx/rest-go-cpp/models"

	"github.com/Dainerx/rest-go-cpp/models/response"
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

	sresponses, err := response.RecentSolveResponses(models.Db, *ptruser, 1)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("RecentSolveResponses() failed %s", err)
		return
	}
	json.NewEncoder(w).Encode(sresponses[0])
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

	limit, err := strconv.ParseInt(r.URL.Query()["count"][0], 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sresponses, err := response.RecentSolveResponses(models.Db, *ptruser, int(limit))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(sresponses)
}
