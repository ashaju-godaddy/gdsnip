package models

import (
	"fmt"
	"net/http"
)

// APIError represents a structured error response
type APIError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// HTTPStatus returns the appropriate HTTP status code for the error
func (e *APIError) HTTPStatus() int {
	switch e.Code {
	case "VALIDATION_ERROR", "MISSING_VARIABLES", "INVALID_SLUG":
		return http.StatusBadRequest
	case "UNAUTHORIZED", "INVALID_TOKEN", "TOKEN_EXPIRED":
		return http.StatusUnauthorized
	case "FORBIDDEN", "INSUFFICIENT_PERMISSIONS":
		return http.StatusForbidden
	case "NOT_FOUND", "SNIPPET_NOT_FOUND", "USER_NOT_FOUND":
		return http.StatusNotFound
	case "CONFLICT", "DUPLICATE_USERNAME", "DUPLICATE_EMAIL", "DUPLICATE_SLUG":
		return http.StatusConflict
	case "RATE_LIMIT_EXCEEDED":
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}

// Error constructors for common error types

// NewValidationError creates a validation error
func NewValidationError(message string, details interface{}) *APIError {
	return &APIError{
		Code:    "VALIDATION_ERROR",
		Message: message,
		Details: details,
	}
}

// NewMissingVariablesError creates an error for missing template variables
func NewMissingVariablesError(message string, details interface{}) *APIError {
	return &APIError{
		Code:    "MISSING_VARIABLES",
		Message: message,
		Details: details,
	}
}

// NewNotFoundError creates a not found error
func NewNotFoundError(resource string) *APIError {
	return &APIError{
		Code:    "NOT_FOUND",
		Message: fmt.Sprintf("%s not found", resource),
	}
}

// NewSnippetNotFoundError creates a snippet not found error
func NewSnippetNotFoundError(path string) *APIError {
	return &APIError{
		Code:    "SNIPPET_NOT_FOUND",
		Message: fmt.Sprintf("snippet '%s' not found", path),
	}
}

// NewUserNotFoundError creates a user not found error
func NewUserNotFoundError(identifier string) *APIError {
	return &APIError{
		Code:    "USER_NOT_FOUND",
		Message: fmt.Sprintf("user '%s' not found", identifier),
	}
}

// NewUnauthorizedError creates an unauthorized error
func NewUnauthorizedError(message string) *APIError {
	if message == "" {
		message = "unauthorized"
	}
	return &APIError{
		Code:    "UNAUTHORIZED",
		Message: message,
	}
}

// NewForbiddenError creates a forbidden error
func NewForbiddenError(message string) *APIError {
	if message == "" {
		message = "insufficient permissions"
	}
	return &APIError{
		Code:    "FORBIDDEN",
		Message: message,
	}
}

// NewConflictError creates a conflict error
func NewConflictError(resource string, details interface{}) *APIError {
	return &APIError{
		Code:    "CONFLICT",
		Message: fmt.Sprintf("%s already exists", resource),
		Details: details,
	}
}

// NewDuplicateUsernameError creates a duplicate username error
func NewDuplicateUsernameError(username string) *APIError {
	return &APIError{
		Code:    "DUPLICATE_USERNAME",
		Message: fmt.Sprintf("username '%s' is already taken", username),
	}
}

// NewDuplicateEmailError creates a duplicate email error
func NewDuplicateEmailError(email string) *APIError {
	return &APIError{
		Code:    "DUPLICATE_EMAIL",
		Message: fmt.Sprintf("email '%s' is already registered", email),
	}
}

// NewDuplicateSlugError creates a duplicate slug error
func NewDuplicateSlugError(namespace, slug string) *APIError {
	return &APIError{
		Code:    "DUPLICATE_SLUG",
		Message: fmt.Sprintf("snippet '%s/%s' already exists", namespace, slug),
	}
}

// NewRateLimitError creates a rate limit exceeded error
func NewRateLimitError() *APIError {
	return &APIError{
		Code:    "RATE_LIMIT_EXCEEDED",
		Message: "too many requests, please try again later",
	}
}

// NewInternalError creates an internal server error
func NewInternalError(message string) *APIError {
	if message == "" {
		message = "internal server error"
	}
	return &APIError{
		Code:    "INTERNAL_ERROR",
		Message: message,
	}
}

// NewInvalidTokenError creates an invalid token error
func NewInvalidTokenError() *APIError {
	return &APIError{
		Code:    "INVALID_TOKEN",
		Message: "invalid or malformed token",
	}
}

// NewTokenExpiredError creates a token expired error
func NewTokenExpiredError() *APIError {
	return &APIError{
		Code:    "TOKEN_EXPIRED",
		Message: "token has expired",
	}
}
