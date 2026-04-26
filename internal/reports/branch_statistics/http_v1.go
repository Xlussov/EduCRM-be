package branch_statistics

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

// Handle godoc
// @Summary      Get Branch Statistics
// @Description  Get statistics for a specific branch including active students, completed/cancelled lessons, and attendance percentage.
// @Tags         reports
// @Accept       json
// @Produce      json
// @Param        branch_id   query     string  true  "Branch ID"
// @Param        start_date  query     string  false "Start Date (YYYY-MM-DD)"
// @Param        end_date    query     string  false "End Date (YYYY-MM-DD)"
// @Success      200         {object}  Response
// @Failure      400         {object}  response.ErrorResponse
// @Failure      401         {object}  response.ErrorResponse
// @Failure      403         {object}  response.ErrorResponse
// @Failure      500         {object}  response.ErrorResponse
// @Security     Bearer
// @Router       /api/v1/reports/branch-statistics [get]
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
		return response.Error(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
	}

	ctx := c.Request().Context()
	res, err := h.usecase.Execute(ctx, *caller, req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrBranchAccessDenied):
			return response.Error(c, http.StatusForbidden, "BRANCH_ACCESS_DENIED", err.Error(), nil)
		default:
			return response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		}
	}

	return response.Success(c, http.StatusOK, res)
}
