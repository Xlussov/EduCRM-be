package refresh

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

func (h *Handler) Handle(c echo.Context) error {
	var req Request
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body", nil)
	}

	if req.RefreshToken == "" {
		return response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Refresh token is required", nil)
	}

	res, err := h.usecase.Execute(c.Request().Context(), req)
	if err != nil {
		if err.Error() == "invalid refresh token" || err.Error() == "refresh token expired" || err.Error() == "token reused" {
			return response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error(), nil)
		}
		if err.Error() == "invalid user" {
			return response.Error(c, http.StatusForbidden, "FORBIDDEN", "User not found", nil)
		}
		return response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error", nil)
	}

	return response.Success(c, http.StatusOK, res)
}
