package main

import (
	"net/http"

	database "github.com/4Noyis/my-library/internal/database"
	"github.com/4Noyis/my-library/internal/handlers"
	"github.com/4Noyis/my-library/internal/logger"
	"github.com/4Noyis/my-library/internal/middleware"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func main() {

	logger.LogInfo("Starting library management server", nil)

	err := database.ConnectMongoDB()
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"type":  "startup",
		}).Fatal("Failed to connect to MongoDB")
	}
	defer database.DisconnectMongoDB()

	// Single handler that routes internally by HTTP method
	// http.HandleFunc("/api/v1/books", handlers.BooksHandler)

	r := mux.NewRouter()

	// Add logging middleware
	r.Use(middleware.LoggingMiddleware)

	// RESTful routes
	r.HandleFunc("/api/v1/books", handlers.GetAllBooksHandler).Methods("GET")
	r.HandleFunc("/api/v1/books", handlers.CreateBookHandler).Methods("POST")
	r.HandleFunc("/api/v1/books/{id}", handlers.GetOneBookHandler).Methods("GET")
	r.HandleFunc("/api/v1/books/{id}", handlers.UpdateBookHandler).Methods("PATCH")
	r.HandleFunc("/api/v1/books/{id}", handlers.DeleteBookHandler).Methods("DELETE")

	logger.LogInfo("Server starting on port 8080", logrus.Fields{"port": 8080})

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"port":  8080,
			"type":  "startup",
		}).Fatal("Server failed to start")
	}

}
