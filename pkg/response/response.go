package response

import (
	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

func Error(c echo.Context, status int, code, message string, details map[string]string) error {
	return c.JSON(status, ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

func Success(c echo.Context, status int, data interface{}) error {
	return c.JSON(status, data)
}

func Message(c echo.Context, status int, message string) error {
	return c.JSON(status, map[string]string{"message": message})
}
