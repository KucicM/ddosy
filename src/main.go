package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	port := flag.Int("port", 4000, "port for the server")

	a := &App{}

	r := mux.NewRouter()
	r.PathPrefix("/schedule").HandlerFunc(a.ScheduelHandler).Methods("POST")
	r.PathPrefix("/status").HandlerFunc(a.StatusHandler).Methods("GET")
	r.PathPrefix("/kill").HandlerFunc(a.KillHandler)

	srv := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf(":%d", *port),
	}

	log.Fatalln(srv.ListenAndServe())
}
