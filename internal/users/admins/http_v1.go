package admins

import (
	"errors"
	"net/http"

	"github.com/Xlussov/EduCRM-be/internal/controller/http/middleware"
	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/Xlussov/EduCRM-be/pkg/response"
	"github.com/Xlussov/EduCRM-be/pkg/validator"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	usecase *UseCase
}

func NewHandler(uc *UseCase) *Handler {
	return &Handler{usecase: uc}
}

// @Summary Create Admin
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body Request true "Admin details"
// @Success 201 {object} Response "Created"
// @Failure 400 {object} response.ErrorResponse "Bad Request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/v1/users/admins [post]
func (h *Handler) Handle(c echo.Context) error {
	userToken, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Missing token", nil)
	}
	userClaims, ok := userToken.Claims.(*middleware.CustomClaims)
	if !ok || userClaims.Role != "SUPERADMIN" {
		return response.Error(c, http.StatusForbidden, "ROLE_NOT_ALLOWED", "Only SUPERADMIN can create ADMINs", nil)
	}

	var req Request
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body", nil)
	}

	if err := c.Validate(&req); err != nil {
		valErrs := validator.ParseError(err)
		return response.Error(c, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid request data", valErrs)
	}

	res, err := h.usecase.Execute(c.Request().Context(), req)
	if err != nil {
		if errors.Is(err, domain.ErrPhoneAlreadyExists) {
			return response.Error(c, http.StatusConflict, "PHONE_ALREADY_EXISTS", "User with this phone already exists", nil)
		}
		return response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
	}

	return response.Success(c, http.StatusCreated, res)
}
