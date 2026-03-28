package login

import (
	"net/http"

	"github.com/Xlussov/EduCRM-be/pkg/response"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	usecase *UseCase
}

func NewHandler(uc *UseCase) *Handler {
	return &Handler{usecase: uc}
}

// @Summary Login
// @Tags auth
// @Accept json
// @Produce json
// @Param request body Request true "Login credentials"
// @Success 200 {object} Response "Success"
// @Failure 400 {object} response.ErrorResponse "Bad Request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/v1/auth/login [post]
func (h *Handler) Handle(c echo.Context) error {
	var req Request
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body", nil)
	}

	if req.Phone == "" || req.Password == "" {
		return response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Phone and password are required", nil)
	}

	res, err := h.usecase.Execute(c.Request().Context(), req)
	if err != nil {
		if err.Error() == "invalid credentials" {
			return response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid phone or password", nil)
		}
		if err.Error() == "user is not active" {
			return response.Error(c, http.StatusForbidden, "FORBIDDEN", "User account is disabled", nil)
		}
		return response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error", nil)
	}

	return response.Success(c, http.StatusOK, res)
}
