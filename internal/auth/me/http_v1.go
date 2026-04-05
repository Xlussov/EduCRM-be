package me

import (
	"net/http"

	"github.com/Xlussov/EduCRM-be/internal/controller/http/middleware"
	"github.com/Xlussov/EduCRM-be/pkg/response"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	usecase *UseCase
}

func NewHandler(uc *UseCase) *Handler {
	return &Handler{
		usecase: uc,
	}
}

// @Summary Get current user profile
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} Response
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/auth/me [get]
func (h *Handler) Handle(c echo.Context) error {
	userToken, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Missing token", nil)
	}

	userClaims, ok := userToken.Claims.(*middleware.CustomClaims)
	if !ok {
		return response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid claims", nil)
	}

	userID, err := uuid.Parse(userClaims.UserID)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid user ID in token", nil)
	}

	res, err := h.usecase.Execute(c.Request().Context(), userID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return response.Error(c, http.StatusNotFound, "NOT_FOUND", "User not found", nil)
		}
		return response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get user data", nil)
	}

	return response.Success(c, http.StatusOK, res)
}
