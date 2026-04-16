package list

import (
	"net/http"

	"github.com/Xlussov/EduCRM-be/internal/controller/http/middleware"
	"github.com/Xlussov/EduCRM-be/pkg/response"
	"github.com/labstack/echo/v4"
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
// @Success 200 {array} TeacherResponse "List of teachers"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/v1/users/teachers [get]
func (h *Handler) Handle(c echo.Context) error {
	caller, err := middleware.GetCaller(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error(), nil)
	}

	res, err := h.usecase.Execute(c.Request().Context(), *caller, Request{})
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
	}

	return response.Success(c, http.StatusOK, res)
}
