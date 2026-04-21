package logout

import (
	"net/http"

	"github.com/Xlussov/EduCRM-be/pkg/response"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	usecase *UseCase
}

func NewHandler(uc *UseCase) *Handler {
	return &Handler{
		usecase: uc,
	}
}

// @Summary Logout user
// @Description Revokes the provided refresh token
// @Tags auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body Request true "Refresh token to revoke"
// @Success 200 {object} Response
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/auth/logout [post]
func (h *Handler) Handle(c echo.Context) error {
	var req Request

	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
	}

	if err := c.Validate(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", nil)
	}

	res, err := h.usecase.Execute(c.Request().Context(), req)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to process logout", nil)
	}

	return response.Success(c, http.StatusOK, res)
}
