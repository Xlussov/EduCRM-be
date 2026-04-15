package get

import (
	"errors"
	"net/http"

	"github.com/Xlussov/EduCRM-be/internal/controller/http/middleware"
	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/Xlussov/EduCRM-be/pkg/response"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	usecase *UseCase
}

func NewHandler(uc *UseCase) *Handler {
	return &Handler{usecase: uc}
}

// @Summary Get Subject by ID
// @Tags subjects
// @Security BearerAuth
// @Produce json
// @Param id path string true "Subject ID" format(uuid)
// @Success 200 {object} Response "Subject details"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden"
// @Failure 404 {object} response.ErrorResponse "Not Found"
// @Failure 400 {object} response.ErrorResponse "Bad Request"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/v1/subjects/{id} [get]
func (h *Handler) Handle(c echo.Context) error {
	caller, err := middleware.GetCaller(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error(), nil)
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid subject id", nil)
	}

	res, err := h.usecase.Execute(c.Request().Context(), *caller, Request{ID: id})
	if err != nil {
		if errors.Is(err, domain.ErrBranchAccessDenied) {
			return response.Error(c, http.StatusForbidden, "BRANCH_ACCESS_DENIED", err.Error(), nil)
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return response.Error(c, http.StatusNotFound, "NOT_FOUND", "Subject not found", nil)
		}
		return response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get subject", nil)
	}
	return response.Success(c, http.StatusOK, res)
}
