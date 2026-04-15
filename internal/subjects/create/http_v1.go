package create

import (
	"errors"
	"net/http"

	"github.com/Xlussov/EduCRM-be/internal/controller/http/middleware"
	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/Xlussov/EduCRM-be/pkg/response"
	"github.com/Xlussov/EduCRM-be/pkg/validator"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	usecase *UseCase
}

func NewHandler(uc *UseCase) *Handler {
	return &Handler{usecase: uc}
}

// @Summary Create Subject
// @Tags subjects
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body Request true "Subject details"
// @Success 201 {object} Response "Created subject"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden"
// @Failure 400 {object} response.ErrorResponse "Bad Request"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/v1/subjects [post]
func (h *Handler) Handle(c echo.Context) error {
	caller, err := middleware.GetCaller(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error(), nil)
	}

	var req Request
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body", nil)
	}

	if err := c.Validate(&req); err != nil {
		valErrs := validator.ParseError(err)
		return response.Error(c, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid request data", valErrs)
	}

	res, err := h.usecase.Execute(c.Request().Context(), *caller, req)
	if err != nil {
		if errors.Is(err, domain.ErrBranchAccessDenied) {
			return response.Error(c, http.StatusForbidden, "BRANCH_ACCESS_DENIED", err.Error(), nil)
		}
		if errors.Is(err, domain.ErrArchivedReference) {
			return response.Error(c, http.StatusBadRequest, "ARCHIVED_REFERENCE", err.Error(), nil)
		}
		if errors.Is(err, domain.ErrAlreadyExists) {
			return response.Error(c, http.StatusConflict, "SUBJECT_ALREADY_EXISTS", "Subject name already exists in this branch", nil)
		}
		return response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create subject", nil)
	}

	return response.Success(c, http.StatusCreated, res)
}
