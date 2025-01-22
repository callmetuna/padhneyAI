package main

import (
	"log"
	"net/http"
	"rag-go-app/api"
)

func main() {
	// Initialize and run the server
	log.Println("Starting server on port 8080...")
	err := http.ListenAndServe(":8080", api.SetupRoutes())
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
