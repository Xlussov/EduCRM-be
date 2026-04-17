package create_template

import (
	"errors"
	"net/http"

	"github.com/Xlussov/EduCRM-be/internal/controller/http/middleware"
	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/Xlussov/EduCRM-be/pkg/response"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	usecase *UseCase
}

func NewHandler(uc *UseCase) *Handler {
	return &Handler{usecase: uc}
}

// @Summary      Create lesson template
// @Description  Creates a new recurring lesson template and schedules lessons
// @Tags         lessons
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body Request true "Template info"
// @Success      201  {object}  Response
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
// @Failure      403  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /api/v1/lessons/templates [post]
func (h *Handler) Handle(c echo.Context) error {
	caller, err := middleware.GetCaller(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error(), nil)
	}

	var req Request
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
	}

	if err := c.Validate(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "VALIDATION_FAILED", err.Error(), nil)
	}

	res, err := h.usecase.Execute(c.Request().Context(), *caller, req)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidInput) {
			return response.Error(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
		}
		if errors.Is(err, domain.ErrBranchAccessDenied) {
			return response.Error(c, http.StatusForbidden, "FORBIDDEN", err.Error(), nil)
		}
		return response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
	}

	return response.Success(c, http.StatusCreated, res)
}
