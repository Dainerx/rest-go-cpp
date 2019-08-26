package main

import (
	"net/http"
	"os"

	"github.com/Dainerx/rest-go-cpp/handlers"
	"github.com/Dainerx/rest-go-cpp/internal"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func init() {
	//File for logging
	f, err := os.OpenFile("var/log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(f)
	// Only log the trace severity or above.
	log.SetLevel(log.TraceLevel)
}
func main() {
	internal.InitDB()
	router := mux.NewRouter()
	router.HandleFunc("/authenticate", handlers.Auth).Methods("POST")
	router.HandleFunc("/solve", handlers.Solve).Methods("POST")
	router.HandleFunc("/status", handlers.Status).Methods("GET")
	router.HandleFunc("/last", handlers.Last).Methods("GET")
	router.HandleFunc("/recent", handlers.Recent).Methods("GET")
	err := http.ListenAndServe(":8000", router)
	if err != nil {
		log.Panicf("Failed to start service: %v.", err)
	}
}
