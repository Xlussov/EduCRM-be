package archive

import (
	"errors"
	"net/http"

	"github.com/Xlussov/EduCRM-be/internal/controller/http/middleware"
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
	return &Handler{
		usecase: uc,
	}
}

// @Summary Update plan status
// @Tags plans
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Plan ID" format(uuid)
// @Param body body Request true "Archived"
// @Success 200 {object} Response
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/plans/{id}/archive [patch]
func (h *Handler) Handle(c echo.Context) error {
	userToken, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Missing token", nil)
	}
	userClaims, ok := userToken.Claims.(*middleware.CustomClaims)
	if !ok {
		return response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid claims", nil)
	}

	if userClaims.Role != "SUPERADMIN" && userClaims.Role != "ADMIN" {
		return response.Error(c, http.StatusForbidden, "ROLE_NOT_ALLOWED", "Only SUPERADMIN or ADMIN can update plans", nil)
	}

	userID, err := uuid.Parse(userClaims.UserID)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid user ID in token", nil)
	}

	idStr := c.Param("id")
	planID, err := uuid.Parse(idStr)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "invalid plan id", nil)
	}

	var req Request
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "invalid request format", nil)
	}

	if err := c.Validate(&req); err != nil {
		valErrs := validator.ParseError(err)
		return response.Error(c, http.StatusBadRequest, "VALIDATION_FAILED", "invalid request data", valErrs)
	}

	res, err := h.usecase.Execute(c.Request().Context(), userID, userClaims.Role, planID, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrBranchAccessDenied):
			return response.Error(c, http.StatusForbidden, "BRANCH_ACCESS_DENIED", err.Error(), nil)
		default:
			if err.Error() == "no rows in result set" {
				return response.Error(c, http.StatusNotFound, "NOT_FOUND", "Plan not found", nil)
			}
			return response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		}
	}

	return response.Success(c, http.StatusOK, res)
}
