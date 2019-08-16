package main

import (
	"net/http"

	"github.com/Dainerx/rest-go-cpp/handlers"
	"github.com/Dainerx/rest-go-cpp/models"
	"github.com/gorilla/mux"
)

func main() {
	models.InitDB()
	router := mux.NewRouter()
	router.HandleFunc("/authenticate", handlers.Auth).Methods("POST")
	router.HandleFunc("/solve", handlers.Solve).Methods("POST")
	router.HandleFunc("/status", handlers.Status).Methods("GET")
	http.ListenAndServe(":8000", router)
}
