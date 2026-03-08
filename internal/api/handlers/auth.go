package handlers

import (
	"net/http"

	"github.com/ashaju-godaddy/gdsnip/internal/api/service"
	"github.com/labstack/echo/v4"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// RegisterRequest represents registration request body
type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginRequest represents login request body
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	User  interface{} `json:"user"`
	Token string      `json:"token"`
}

// Register handles user registration
func (h *AuthHandler) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return Error(c, err)
	}

	user, token, err := h.authService.Register(req.Email, req.Username, req.Password)
	if err != nil {
		return Error(c, err)
	}

	// Don't return password hash to client
	userResponse := map[string]interface{}{
		"id":         user.ID,
		"email":      user.Email,
		"username":   user.Username,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	}

	return Success(c, http.StatusCreated, AuthResponse{
		User:  userResponse,
		Token: token,
	})
}

// Login handles user login
func (h *AuthHandler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return Error(c, err)
	}

	user, token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		return Error(c, err)
	}

	// Don't return password hash to client
	userResponse := map[string]interface{}{
		"id":         user.ID,
		"email":      user.Email,
		"username":   user.Username,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	}

	return Success(c, http.StatusOK, AuthResponse{
		User:  userResponse,
		Token: token,
	})
}

// Me returns current authenticated user
func (h *AuthHandler) Me(c echo.Context) error {
	userID := c.Get("user_id").(string)

	user, err := h.authService.GetCurrentUser(userID)
	if err != nil {
		return Error(c, err)
	}

	// Don't return password hash to client
	userResponse := map[string]interface{}{
		"id":         user.ID,
		"email":      user.Email,
		"username":   user.Username,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	}

	return Success(c, http.StatusOK, userResponse)
}
