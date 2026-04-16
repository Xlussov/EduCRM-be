package archive

import (
	"errors"
	"net/http"

	"github.com/Xlussov/EduCRM-be/internal/controller/http/middleware"
	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/Xlussov/EduCRM-be/pkg/response"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	usecase *UseCase
}

func NewHandler(uc *UseCase) *Handler {
	return &Handler{usecase: uc}
}

// @Summary Archive Teacher
// @Tags teachers
// @Security BearerAuth
// @Produce json
// @Param id path string true "Teacher ID" format(uuid)
// @Success 200 {object} Response "Archived"
// @Failure 400 {object} response.ErrorResponse "Bad Request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden"
// @Failure 404 {object} response.ErrorResponse "Not Found"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/v1/users/teachers/{id}/archive [patch]
func (h *Handler) Handle(c echo.Context) error {
	caller, err := middleware.GetCaller(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error(), nil)
	}

	teacherID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid teacher ID", nil)
	}

	res, err := h.usecase.Execute(c.Request().Context(), *caller, teacherID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrBranchAccessDenied):
			return response.Error(c, http.StatusForbidden, "BRANCH_ACCESS_DENIED", err.Error(), nil)
		case errors.Is(err, domain.ErrAlreadyArchived):
			return response.Error(c, http.StatusBadRequest, "ALREADY_ARCHIVED", "This teacher is already archived", nil)
		case errors.Is(err, domain.ErrNotFound):
			return response.Error(c, http.StatusNotFound, "NOT_FOUND", "Teacher not found", nil)
		default:
			return response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		}
	}

	return response.Success(c, http.StatusOK, res)
}
