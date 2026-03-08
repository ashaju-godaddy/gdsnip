package handlers

import (
	"github.com/ashaju-godaddy/gdsnip/pkg/models"
	"github.com/labstack/echo/v4"
)

// Success returns a successful JSON response
func Success(c echo.Context, status int, data interface{}) error {
	return c.JSON(status, map[string]interface{}{
		"success": true,
		"data":    data,
	})
}

// Error returns an error JSON response
func Error(c echo.Context, err error) error {
	// Check if it's an APIError
	if apiErr, ok := err.(*models.APIError); ok {
		return c.JSON(apiErr.HTTPStatus(), map[string]interface{}{
			"success": false,
			"error":   apiErr,
		})
	}

	// Default internal server error
	apiErr := models.NewInternalError(err.Error())
	return c.JSON(apiErr.HTTPStatus(), map[string]interface{}{
		"success": false,
		"error":   apiErr,
	})
}

// PaginatedSuccess returns a paginated successful JSON response
func PaginatedSuccess(c echo.Context, status int, data interface{}, total, limit, offset int) error {
	return c.JSON(status, map[string]interface{}{
		"success": true,
		"data":    data,
		"pagination": map[string]interface{}{
			"total":  total,
			"limit":  limit,
			"offset": offset,
		},
	})
}
