package update

import (
	"errors"
	"net/http"

	"github.com/Xlussov/EduCRM-be/internal/controller/http/middleware"
	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/Xlussov/EduCRM-be/pkg/response"
	"github.com/Xlussov/EduCRM-be/pkg/validator"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	usecase *UseCase
}

func NewHandler(uc *UseCase) *Handler {
	return &Handler{usecase: uc}
}

// @Summary Update Admin
// @Tags admins
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Admin ID" format(uuid)
// @Param request body Request true "Admin update payload"
// @Success 200 {object} Response "Updated admin"
// @Failure 400 {object} response.ErrorResponse "Bad Request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden"
// @Failure 404 {object} response.ErrorResponse "Not Found"
// @Failure 409 {object} response.ErrorResponse "Conflict"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/v1/users/admins/{id} [put]
func (h *Handler) Handle(c echo.Context) error {
	userToken, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Missing token", nil)
	}
	userClaims, ok := userToken.Claims.(*middleware.CustomClaims)
	if !ok || userClaims.Role != "SUPERADMIN" {
		return response.Error(c, http.StatusForbidden, "ROLE_NOT_ALLOWED", "Only SUPERADMIN can update ADMINs", nil)
	}

	adminID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid admin ID", nil)
	}

	var req Request
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body", nil)
	}
	if err := c.Validate(&req); err != nil {
		valErrs := validator.ParseError(err)
		return response.Error(c, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid request data", valErrs)
	}

	res, err := h.usecase.Execute(c.Request().Context(), adminID, req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrPhoneAlreadyExists):
			return response.Error(c, http.StatusConflict, "PHONE_ALREADY_EXISTS", "User with this phone already exists", nil)
		case errors.Is(err, domain.ErrArchivedReference):
			return response.Error(c, http.StatusBadRequest, "ARCHIVED_REFERENCE", "Cannot reference archived entity", nil)
		case errors.Is(err, domain.ErrCannotEditArchived):
			return response.Error(c, http.StatusBadRequest, "CANNOT_EDIT_ARCHIVED", err.Error(), nil)
		case errors.Is(err, domain.ErrNotFound):
			return response.Error(c, http.StatusNotFound, "NOT_FOUND", "Admin not found", nil)
		default:
			return response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		}
	}

	return response.Success(c, http.StatusOK, res)
}
