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

// @Summary Archive Subject
// @Tags subjects
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Subject ID format(uuid)"
// @Success 200 {object} map[string]string "Success message"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden"
// @Failure 400 {object} response.ErrorResponse "Bad Request"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/v1/subjects/{id}/archive [patch]
func (h *Handler) Handle(c echo.Context) error {
	caller, err := middleware.GetCaller(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error(), nil)
	}

	subjectIDParam := c.Param("id")
	subjectID, err := uuid.Parse(subjectIDParam)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid subject ID", nil)
	}

	res, err := h.usecase.Execute(c.Request().Context(), *caller, subjectID)
	if err != nil {
		if errors.Is(err, domain.ErrBranchAccessDenied) {
			return response.Error(c, http.StatusForbidden, "BRANCH_ACCESS_DENIED", err.Error(), nil)
		}
		if errors.Is(err, domain.ErrAlreadyArchived) {
			return response.Error(c, http.StatusBadRequest, "ALREADY_ARCHIVED", "This subject is already in the archive", nil)
		}
		return response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to archive subject", nil)
	}

	return response.Success(c, http.StatusNoContent, res)
}
