package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/4Noyis/my-library/internal/logger"
	"github.com/4Noyis/my-library/internal/models"
	"github.com/4Noyis/my-library/internal/services"
	"github.com/sirupsen/logrus"
)

var userService = services.NewUserService()

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"type":  "validation",
		}).Error("Invalid request body for registration")

		response := models.UserResponse{
			Status:  "error",
			Message: "Invalid request body",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	user, err := userService.RegisterUser(&req)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"error":    err.Error(),
			"username": req.Username,
			"email":    req.Email,
			"type":     "registration",
		}).Error("User registration failed")

		var statusCode int
		if err.Error() == "username already exists" || err.Error() == "email already exists" {
			statusCode = http.StatusConflict
		} else {
			statusCode = http.StatusInternalServerError
		}

		response := models.UserResponse{
			Status:  "error",
			Message: err.Error(),
		}
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(response)
		return
	}

	logger.Logger.WithFields(logrus.Fields{
		"user_id":  user.ID.Hex(),
		"username": user.Username,
		"email":    user.Email,
		"type":     "registration",
	}).Info("User registered successfully")

	response := models.UserResponse{
		Status:  "success",
		Message: "User registered successfully",
		Data:    user,
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"type":  "validation",
		}).Error("Invalid request body for login")

		response := models.UserResponse{
			Status:  "error",
			Message: "Invalid request body",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	loginResponse, err := userService.LoginUser(&req)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"error":    err.Error(),
			"username": req.Username,
			"type":     "login",
		}).Error("User login failed")

		var statusCode int
		if err.Error() == "invalid credentials" || err.Error() == "account is deactivated" {
			statusCode = http.StatusUnauthorized
		} else {
			statusCode = http.StatusInternalServerError
		}

		response := models.UserResponse{
			Status:  "error",
			Message: err.Error(),
		}
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(response)
		return
	}

	logger.Logger.WithFields(logrus.Fields{
		"user_id":  loginResponse.User.ID.Hex(),
		"username": loginResponse.User.Username,
		"type":     "login",
	}).Info("User logged in successfully")

	response := models.UserResponse{
		Status:  "success",
		Message: "Login successful",
		Data:    loginResponse,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
