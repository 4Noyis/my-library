package services

import (
	"errors"
	"os"
	"time"

	"github.com/4Noyis/my-library/internal/models"
	"github.com/4Noyis/my-library/internal/repositories"
	"github.com/golang-jwt/jwt/v5"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo  *repositories.UserRepository
	jwtSecret []byte
}

func NewUserService() *UserService {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-super-secret-jwt-key-change-in-production" // Default for development
	}

	return &UserService{
		userRepo:  repositories.NewUserRepository(),
		jwtSecret: []byte(secret),
	}
}

func (us *UserService) RegisterUser(req *models.RegisterRequest) (*models.User, error) {
	// Check if username already exists
	existingUser, err := us.userRepo.GetUserByUsername(req.Username)
	if err != nil && err != mongo.ErrNoDocuments {
		// Database error occurred
		return nil, errors.New("database error while checking username")
	}
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// Check if email already exists
	existingUser, err = us.userRepo.GetUserByEmail(req.Email)
	if err != nil && err != mongo.ErrNoDocuments {
		// Database error occurred
		return nil, errors.New("database error while checking email")
	}
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Set default role if not provided
	role := req.Role
	if role == "" {
		role = "user"
	}

	// Create user
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     role,
	}

	err = us.userRepo.CreateUser(user)
	if err != nil {
		return nil, errors.New("failed to create user")
	}

	// Don't return password
	user.Password = ""
	return user, nil
}

func (us *UserService) LoginUser(req *models.LoginRequest) (*models.LoginResponse, error) {
	// Get user by username
	user, err := us.userRepo.GetUserByUsername(req.Username)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("invalid credentials")
		}
		return nil, errors.New("database error during login")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid credentials (password)")
	}

	// Generate JWT token
	token, err := us.generateJWT(user)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	// Don't return password
	user.Password = ""

	return &models.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (us *UserService) generateJWT(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID.Hex(),
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(us.jwtSecret)
}

func (us *UserService) ValidateJWT(tokenString string) (*models.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return us.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			return nil, errors.New("invalid user_id in token")
		}

		userID, err := primitive.ObjectIDFromHex(userIDStr)
		if err != nil {
			return nil, errors.New("invalid user_id format")
		}

		user, err := us.userRepo.GetUserByID(userID)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, errors.New("user not found")
			}
			return nil, errors.New("database error during token validation")
		}

		if !user.IsActive {
			return nil, errors.New("user account is deactivated")
		}

		return user, nil
	}

	return nil, errors.New("invalid token")
}
