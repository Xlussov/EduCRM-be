package list

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

// @Summary List Students
// @Tags students
// @Security BearerAuth
// @Produce json
// @Param branch_id query string true "Branch ID" format(uuid)
// @Param search query string false "Search by first or last name"
// @Param status query string false "Filter by status" Enums(ACTIVE, ARCHIVED)
// @Success 200 {object} Response "List of students"
// @Failure 400 {object} response.ErrorResponse "Bad Request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/v1/students [get]
func (h *Handler) Handle(c echo.Context) error {
	caller, err := middleware.GetCaller(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error(), nil)
	}

	branchIDStr := c.QueryParam("branch_id")
	branchID, err := uuid.Parse(branchIDStr)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid or missing branch_id", nil)
	}

	req := Request{
		BranchID: branchID,
		Search:   c.QueryParam("search"),
		Status:   c.QueryParam("status"),
	}

	res, err := h.usecase.Execute(c.Request().Context(), *caller, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrBranchIDRequired):
			return response.Error(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
		case errors.Is(err, domain.ErrInvalidInput):
			return response.Error(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
		case errors.Is(err, domain.ErrBranchAccessDenied):
			return response.Error(c, http.StatusForbidden, "BRANCH_ACCESS_DENIED", err.Error(), nil)
		default:
			return response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		}
	}

	return response.Success(c, http.StatusOK, res)
}
