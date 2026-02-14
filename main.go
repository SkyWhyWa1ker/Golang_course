package main

import (
	"log"
	"net/http"

	"practice2/handlers"
	"practice2/middleware"
)

func main() {

	mux := http.NewServeMux()

	// endpoint
	mux.HandleFunc("/tasks", handlers.TasksHandler)

	// middleware chain
	handler := middleware.Logging(
		middleware.APIKey(mux),
	)

	log.Println("Server started on :8080")

	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatal(err)
	}
}
