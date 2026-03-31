package update

import (
	"net/http"

	"github.com/Xlussov/EduCRM-be/internal/controller/http/middleware"
	"github.com/Xlussov/EduCRM-be/pkg/response"
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

// @Summary Update Branch
// @Tags branches
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Branch ID format(uuid)"
// @Param request body Request true "Update branch details"
// @Success 200 {object} map[string]string "Success message"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden"
// @Failure 400 {object} response.ErrorResponse "Bad Request"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/v1/branches/{id} [put]
func (h *Handler) Handle(c echo.Context) error {
	userToken, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Missing token", nil)
	}
	userClaims, ok := userToken.Claims.(*middleware.CustomClaims)
	if !ok {
		return response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid claims", nil)
	}

	branchIDParam := c.Param("id")
	branchID, err := uuid.Parse(branchIDParam)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid branch ID", nil)
	}

	if userClaims.Role != "SUPERADMIN" && userClaims.Role != "ADMIN" {
		return response.Error(c, http.StatusForbidden, "ROLE_NOT_ALLOWED", "Only SUPERADMIN or ADMIN can update branches", nil)
	}

	if userClaims.Role == "ADMIN" {
		hasAccess := false
		for _, b := range userClaims.BranchIDs {
			if b == branchID.String() {
				hasAccess = true
				break
			}
		}
		if !hasAccess {
			return response.Error(c, http.StatusForbidden, "BRANCH_ACCESS_DENIED", "Admin cannot access this branch", nil)
		}
	}

	var req Request
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body", nil)
	}

	res, err := h.usecase.Execute(c.Request().Context(), branchID, req)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
	}

	return response.Success(c, http.StatusNoContent, res)
}
