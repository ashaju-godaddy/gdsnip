package middleware

import (
	"strings"

	"github.com/ashaju-godaddy/gdsnip/internal/api/service"
	"github.com/ashaju-godaddy/gdsnip/pkg/models"
	"github.com/labstack/echo/v4"
)

// JWTMiddleware creates middleware that validates JWT tokens
func JWTMiddleware(authService *service.AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extract token from Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				apiErr := models.NewUnauthorizedError("missing authorization header")
				return c.JSON(apiErr.HTTPStatus(), map[string]interface{}{
					"success": false,
					"error":   apiErr,
				})
			}

			// Check Bearer prefix
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				apiErr := models.NewUnauthorizedError("invalid authorization header format")
				return c.JSON(apiErr.HTTPStatus(), map[string]interface{}{
					"success": false,
					"error":   apiErr,
				})
			}

			token := parts[1]

			// Parse and validate token
			claims, err := authService.ParseJWT(token)
			if err != nil {
				apiErr := err.(*models.APIError)
				return c.JSON(apiErr.HTTPStatus(), map[string]interface{}{
					"success": false,
					"error":   apiErr,
				})
			}

			// Set user info in context
			c.Set("user_id", claims.UserID)
			c.Set("username", claims.Username)

			return next(c)
		}
	}
}

// OptionalJWTMiddleware creates middleware that extracts JWT if present but doesn't fail if missing
// This is useful for endpoints that work for both authenticated and unauthenticated users
func OptionalJWTMiddleware(authService *service.AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extract token from Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				// No auth header - continue without setting user_id
				return next(c)
			}

			// Check Bearer prefix
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				// Invalid format - continue without setting user_id
				return next(c)
			}

			token := parts[1]

			// Parse and validate token
			claims, err := authService.ParseJWT(token)
			if err != nil {
				// Invalid token - continue without setting user_id
				return next(c)
			}

			// Set user info in context
			c.Set("user_id", claims.UserID)
			c.Set("username", claims.Username)

			return next(c)
		}
	}
}
