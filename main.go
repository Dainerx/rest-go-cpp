package main

import (
	"betell-rest/handlers"
	"betell-rest/models"
	"net/http"

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
