package list

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/Xlussov/EduCRM-be/internal/controller/http/middleware"
	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/Xlussov/EduCRM-be/pkg/response"
)

type Handler struct {
	usecase *UseCase
}

func NewHandler(uc *UseCase) *Handler {
	return &Handler{usecase: uc}
}

// @Summary List Teachers
// @Tags teachers
// @Security BearerAuth
// @Produce json
// @Param branch_id query string false "Branch ID" format(uuid)
// @Success 200 {array} TeacherResponse "List of teachers"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden"
// @Failure 400 {object} response.ErrorResponse "Bad Request"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/v1/users/teachers [get]
func (h *Handler) Handle(c echo.Context) error {
	caller, err := middleware.GetCaller(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error(), nil)
	}

	var branchID *uuid.UUID
	if raw := c.QueryParam("branch_id"); raw != "" {
		id, err := uuid.Parse(raw)
		if err != nil {
			return response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "invalid branch_id", nil)
		}
		branchID = &id
	}

	res, err := h.usecase.Execute(c.Request().Context(), *caller, Request{BranchID: branchID})
	if err != nil {
		if errors.Is(err, domain.ErrInvalidInput) {
			return response.Error(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
		}
		if errors.Is(err, domain.ErrBranchAccessDenied) {
			return response.Error(c, http.StatusForbidden, "FORBIDDEN", err.Error(), nil)
		}
		return response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
	}

	return response.Success(c, http.StatusOK, res)
}
