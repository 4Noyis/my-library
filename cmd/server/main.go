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

	// Public routes (no authentication required)
	r.HandleFunc("/api/v1/auth/register", handlers.RegisterHandler).Methods("POST")
	r.HandleFunc("/api/v1/auth/login", handlers.LoginHandler).Methods("POST")

	// Protected routes (authentication required)
	protected := r.PathPrefix("/api/v1").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	
	// Book routes - all require authentication
	protected.HandleFunc("/books", handlers.GetAllBooksHandler).Methods("GET")
	protected.HandleFunc("/books", handlers.CreateBookHandler).Methods("POST")
	protected.HandleFunc("/books/{id}", handlers.GetOneBookHandler).Methods("GET")
	protected.HandleFunc("/books/{id}", handlers.UpdateBookHandler).Methods("PATCH")
	protected.HandleFunc("/books/{id}", handlers.DeleteBookHandler).Methods("DELETE")

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
