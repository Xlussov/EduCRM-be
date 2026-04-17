package cancel

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

// @Summary      Cancel lesson
// @Description  Cancels a specific lesson by ID
// @Tags         lessons
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "Lesson ID"
// @Success      200  {object}  Response
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
// @Failure      403  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /api/v1/lessons/{id}/cancel [patch]
func (h *Handler) Handle(c echo.Context) error {
	caller, err := middleware.GetCaller(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error(), nil)
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "invalid lesson ID", nil)
	}

	res, err := h.usecase.Execute(c.Request().Context(), *caller, id)
	if err != nil {
		if errors.Is(err, domain.ErrBranchAccessDenied) {
			return response.Error(c, http.StatusForbidden, "FORBIDDEN", err.Error(), nil)
		}
		if errors.Is(err, domain.ErrInvalidInput) {
			return response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "lesson is already cancelled or completed", nil)
		}
		return response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
	}

	return response.Success(c, http.StatusOK, res)
}
