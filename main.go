package main

import (
	"log"
	"net/http"

	"github.com/aryanbroy/rest_go/api"
)

func main() {
	mux := http.NewServeMux()
	server := api.NewTaskServer()

	// mux.HandleFunc("GET /test", api.TestHandler)
	mux.HandleFunc("POST /task/", server.CreateTaskHandler)

	log.Println("Starting server at port :8080...")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

