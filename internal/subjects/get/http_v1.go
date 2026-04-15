package get

import (
	"errors"
	"net/http"

	"github.com/Xlussov/EduCRM-be/internal/controller/http/middleware"
	"github.com/Xlussov/EduCRM-be/pkg/response"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	usecase *UseCase
}

func NewHandler(uc *UseCase) *Handler {
	return &Handler{usecase: uc}
}

// @Summary Get Subject by ID
// @Tags subjects
// @Security BearerAuth
// @Produce json
// @Param id path string true "Subject ID" format(uuid)
// @Success 200 {object} Response "Subject details"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden"
// @Failure 404 {object} response.ErrorResponse "Not Found"
// @Failure 400 {object} response.ErrorResponse "Bad Request"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/v1/subjects/{id} [get]
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
		return response.Error(c, http.StatusForbidden, "ROLE_NOT_ALLOWED", "Only SUPERADMIN or ADMIN can get subjects", nil)
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid subject id", nil)
	}

	res, err := h.usecase.Execute(c.Request().Context(), Request{ID: id})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return response.Error(c, http.StatusNotFound, "NOT_FOUND", "Subject not found", nil)
		}
		return response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get subject", nil)
	}

	if userClaims.Role == "ADMIN" {
		allowed := false
		for _, bID := range userClaims.BranchIDs {
			if bID == res.Subject.BranchID {
				allowed = true
				break
			}
		}
		if !allowed {
			return response.Error(c, http.StatusForbidden, "FORBIDDEN", "You do not have access to this branch", nil)
		}
	}

	return response.Success(c, http.StatusOK, res)
}
