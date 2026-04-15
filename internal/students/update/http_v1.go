package update

import (
	"errors"
	"net/http"

	"github.com/Xlussov/EduCRM-be/internal/controller/http/middleware"
	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/Xlussov/EduCRM-be/pkg/response"
	"github.com/Xlussov/EduCRM-be/pkg/validator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	usecase *UseCase
}

func NewHandler(uc *UseCase) *Handler {
	return &Handler{usecase: uc}
}

// @Summary Update Student
// @Tags students
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Student ID" format(uuid)
// @Param request body Request true "Student details"
// @Success 200 {object} Response "Updated student"
// @Failure 400 {object} response.ErrorResponse "Bad Request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/v1/students/{id} [put]
func (h *Handler) Handle(c echo.Context) error {
	caller, err := middleware.GetCaller(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error(), nil)
	}

	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid student ID", nil)
	}

	var req Request
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body", nil)
	}

	if err := c.Validate(&req); err != nil {
		valErrs := validator.ParseError(err)
		return response.Error(c, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid request data", valErrs)
	}

	res, err := h.usecase.Execute(c.Request().Context(), *caller, studentID, req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrCannotEditArchived):
			return response.Error(c, http.StatusBadRequest, "CANNOT_EDIT_ARCHIVED", err.Error(), nil)
		case errors.Is(err, domain.ErrBranchAccessDenied):
			return response.Error(c, http.StatusForbidden, "BRANCH_ACCESS_DENIED", err.Error(), nil)
		default:
			return response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		}
	}

	return response.Success(c, http.StatusOK, res)
}
