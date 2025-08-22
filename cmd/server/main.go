package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"port":  port,
			"type":  "startup",
		}).Fatal("Server failed to start")
	}

	// server starts in goroutine
	go func() {
		logger.LogInfo("Server starting", logrus.Fields{"port": port})
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Logger.WithFields(logrus.Fields{
				"error": err.Error(),
				"port":  port,
				"type":  "startup",
			}).Fatal("Server failed to start")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	logger.LogInfo("Server ready - waiting for requests", logrus.Fields{"port": port})

	// block until recieve a signal
	<-quit
	logger.LogInfo("Shutdown signal received, initiating graceful shutdown", logrus.Fields{"port": port})

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.LogInfo("Server forced to shutdown", logrus.Fields{"port": port})
	}

	logger.LogInfo("Server exited gracefully", logrus.Fields{"port": port})

}
