package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/ashaju-godaddy/gdsnip/internal/api/repository"
	"github.com/ashaju-godaddy/gdsnip/pkg/models"
	"github.com/ashaju-godaddy/gdsnip/pkg/validator"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles authentication logic
type AuthService struct {
	userRepo  *repository.UserRepo
	jwtSecret string
	jwtExpiry time.Duration
}

// Claims represents JWT claims
type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo *repository.UserRepo, jwtSecret string, jwtExpiry time.Duration) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiry,
	}
}

// Register creates a new user account
func (s *AuthService) Register(email, username, password string) (*models.User, string, error) {
	// Validate inputs
	if err := validator.ValidateEmail(email); err != nil {
		return nil, "", models.NewValidationError(err.Error(), nil)
	}

	if err := validator.ValidateUsername(username); err != nil {
		return nil, "", models.NewValidationError(err.Error(), nil)
	}

	if err := validator.ValidatePassword(password); err != nil {
		return nil, "", models.NewValidationError(err.Error(), nil)
	}

	// Hash password
	passwordHash, err := s.hashPassword(password)
	if err != nil {
		return nil, "", models.NewInternalError("failed to hash password")
	}

	// Create user
	user := &models.User{
		ID:           uuid.New().String(),
		Email:        email,
		Username:     username,
		PasswordHash: passwordHash,
	}

	if err := s.userRepo.Create(user); err != nil {
		// Check if it's already an APIError (like duplicate email/username)
		if apiErr, ok := err.(*models.APIError); ok {
			return nil, "", apiErr
		}
		return nil, "", models.NewInternalError(fmt.Sprintf("failed to create user: %v", err))
	}

	// Generate JWT token
	token, err := s.generateJWT(user.ID, user.Username)
	if err != nil {
		return nil, "", models.NewInternalError("failed to generate token")
	}

	return user, token, nil
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(email, password string) (*models.User, string, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		// Check if it's a not found error
		if apiErr, ok := err.(*models.APIError); ok && apiErr.Code == "USER_NOT_FOUND" {
			return nil, "", models.NewUnauthorizedError("invalid email or password")
		}
		return nil, "", models.NewInternalError(fmt.Sprintf("failed to get user: %v", err))
	}

	// Check password
	if !s.checkPassword(user.PasswordHash, password) {
		return nil, "", models.NewUnauthorizedError("invalid email or password")
	}

	// Generate JWT token
	token, err := s.generateJWT(user.ID, user.Username)
	if err != nil {
		return nil, "", models.NewInternalError("failed to generate token")
	}

	return user, token, nil
}

// GetCurrentUser retrieves the current user by ID
func (s *AuthService) GetCurrentUser(userID string) (*models.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// hashPassword hashes a password using bcrypt
func (s *AuthService) hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// checkPassword checks if a password matches the hash
func (s *AuthService) checkPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// generateJWT generates a JWT token for a user
func (s *AuthService) generateJWT(userID, username string) (string, error) {
	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.jwtExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseJWT parses and validates a JWT token
func (s *AuthService) ParseJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		// Check if error is due to token expiration
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, models.NewTokenExpiredError()
		}
		return nil, models.NewInvalidTokenError()
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, models.NewInvalidTokenError()
	}

	return claims, nil
}
