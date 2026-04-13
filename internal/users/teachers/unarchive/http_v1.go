package unarchive

import (
	"errors"
	"net/http"

	"github.com/Xlussov/EduCRM-be/internal/controller/http/middleware"
	"github.com/Xlussov/EduCRM-be/internal/domain"
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

// @Summary Unarchive Teacher
// @Tags teachers
// @Security BearerAuth
// @Produce json
// @Param id path string true "Teacher ID" format(uuid)
// @Success 200 {object} Response "Unarchived"
// @Failure 400 {object} response.ErrorResponse "Bad Request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden"
// @Failure 404 {object} response.ErrorResponse "Not Found"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/v1/users/teachers/{id}/unarchive [patch]
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
		return response.Error(c, http.StatusForbidden, "ROLE_NOT_ALLOWED", "Only SUPERADMIN or ADMIN can unarchive TEACHERs", nil)
	}

	teacherID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid teacher ID", nil)
	}

	adminBranchIDs := make([]uuid.UUID, 0, len(userClaims.BranchIDs))
	if userClaims.Role == "ADMIN" {
		for _, rawID := range userClaims.BranchIDs {
			id, err := uuid.Parse(rawID)
			if err != nil {
				return response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid claims", nil)
			}
			adminBranchIDs = append(adminBranchIDs, id)
		}
	}

	res, err := h.usecase.Execute(c.Request().Context(), userClaims.Role, adminBranchIDs, teacherID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrBranchAccessDenied):
			return response.Error(c, http.StatusForbidden, "BRANCH_ACCESS_DENIED", err.Error(), nil)
		case errors.Is(err, domain.ErrAlreadyActive):
			return response.Error(c, http.StatusBadRequest, "ALREADY_ACTIVE", "This teacher is already active", nil)
		case errors.Is(err, domain.ErrNotFound):
			return response.Error(c, http.StatusNotFound, "NOT_FOUND", "Teacher not found", nil)
		default:
			return response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		}
	}

	return response.Success(c, http.StatusOK, res)
}
