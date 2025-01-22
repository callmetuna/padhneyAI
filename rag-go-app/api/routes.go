package api

import (
	"net/http"

	"rag-go-app/service"

	"github.com/gorilla/mux"
)

// Routes struct to hold the services
type Routes struct {
	DataService *service.DataService
	FileService *service.FileService
}

// NewRouter initializes the router and sets up the routes
func NewRouter(routes Routes) *mux.Router {
	router := mux.NewRouter()

	// Document routes
	router.HandleFunc("/documents", createDocumentHandler(routes.DataService)).Methods("POST")
	router.HandleFunc("/documents/{id:[0-9]+}", getDocumentHandler(routes.DataService)).Methods("GET")
	router.HandleFunc("/documents/{id:[0-9]+}", updateDocumentHandler(routes.DataService)).Methods("PUT")

	// File routes
	router.HandleFunc("/files", uploadFileHandler(routes.FileService)).Methods("POST")
	router.HandleFunc("/files/{id:[0-9]+}", getFileHandler(routes.FileService)).Methods("GET")
	router.HandleFunc("/files/{id:[0-9]+}", deleteFileHandler(routes.FileService)).Methods("DELETE")

	return router
}

// Handler functions for documents
func createDocumentHandler(dataService *service.DataService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Logic to create a document
	}
}

func getDocumentHandler(dataService *service.DataService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Logic to get a document by ID
	}
}

func updateDocumentHandler(dataService *service.DataService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Logic to update a document by ID
	}
}

// Handler functions for files
func uploadFileHandler(fileService *service.FileService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Logic to upload a file
	}
}

func getFileHandler(fileService *service.FileService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Logic to get a file by ID
	}
}

func deleteFileHandler(fileService *service.FileService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Logic to delete a file by ID
	}
}
