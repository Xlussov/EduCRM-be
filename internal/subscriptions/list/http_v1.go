package list

import (
	"errors"
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
	return &Handler{usecase: uc}
}

// Handle handles the request to list a student's subscriptions.
// @Summary List Student Subscriptions
// @Description Returns all active subscriptions for a student
// @Tags subscriptions
// @Security BearerAuth
// @Produce json
// @Param id path string true "Student ID" format(uuid)
// @Success 200 {array} SubscriptionResponse "List of student subscriptions"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/v1/students/{id}/subscriptions [get]
func (h *Handler) Handle(c echo.Context) error {
	userToken, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Missing token", nil)
	}
	userClaims, ok := userToken.Claims.(*middleware.CustomClaims)
	if !ok {
		return response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid claims", nil)
	}

	if userClaims.Role != "SUPERADMIN" && userClaims.Role != "ADMIN" {
		return response.Error(c, http.StatusForbidden, "ROLE_NOT_ALLOWED", "Only SUPERADMIN or ADMIN can view student subscriptions", nil)
	}

	userID, err := uuid.Parse(userClaims.UserID)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid user ID in token", nil)
	}

	studentIDStr := c.Param("id")
	studentID, err := uuid.Parse(studentIDStr)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid student ID format", nil)
	}

	res, err := h.usecase.Execute(c.Request().Context(), userID, studentID, userClaims.Role)
	if err != nil {
		switch {
		case errors.Is(err, ErrBranchAccessDenied):
			return response.Error(c, http.StatusForbidden, "BRANCH_ACCESS_DENIED", err.Error(), nil)
		default:
			return response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		}
	}

	return response.Success(c, http.StatusOK, res)
}
