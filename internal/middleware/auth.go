package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/4Noyis/my-library/internal/logger"
	"github.com/4Noyis/my-library/internal/models"
	"github.com/4Noyis/my-library/internal/services"
	"github.com/sirupsen/logrus"
)

type contextKey string

const UserContextKey contextKey = "user"

var userService = services.NewUserService()

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logger.Logger.WithFields(logrus.Fields{
				"path":   r.URL.Path,
				"method": r.Method,
				"type":   "auth",
			}).Error("Missing authorization header")

			response := models.Response{
				Status:  "error",
				Message: "Authorization header required",
			}
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Check if the header starts with "Bearer "
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			logger.Logger.WithFields(logrus.Fields{
				"path":   r.URL.Path,
				"method": r.Method,
				"type":   "auth",
			}).Error("Invalid authorization header format")

			response := models.Response{
				Status:  "error",
				Message: "Invalid authorization header format",
			}
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Validate the JWT token
		token := tokenParts[1]
		user, err := userService.ValidateJWT(token)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{
				"path":   r.URL.Path,
				"method": r.Method,
				"error":  err.Error(),
				"type":   "auth",
			}).Error("Token validation failed")

			response := models.Response{
				Status:  "error",
				Message: "Invalid or expired token",
			}
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Add user to request context
		ctx := context.WithValue(r.Context(), UserContextKey, user)
		r = r.WithContext(ctx)

		logger.Logger.WithFields(logrus.Fields{
			"path":     r.URL.Path,
			"method":   r.Method,
			"user_id":  user.ID.Hex(),
			"username": user.Username,
			"role":     user.Role,
			"type":     "auth",
		}).Info("User authenticated successfully")

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// AdminMiddleware ensures only admin users can access the route
func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(UserContextKey).(*models.User)
		if !ok {
			response := models.Response{
				Status:  "error",
				Message: "User context not found",
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		if user.Role != "admin" {
			logger.Logger.WithFields(logrus.Fields{
				"path":     r.URL.Path,
				"method":   r.Method,
				"user_id":  user.ID.Hex(),
				"username": user.Username,
				"role":     user.Role,
				"type":     "authorization",
			}).Error("Access denied: admin role required")

			response := models.Response{
				Status:  "error",
				Message: "Admin access required",
			}
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(response)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// GetUserFromContext extracts the user from the request context
func GetUserFromContext(r *http.Request) (*models.User, bool) {
	user, ok := r.Context().Value(UserContextKey).(*models.User)
	return user, ok
}