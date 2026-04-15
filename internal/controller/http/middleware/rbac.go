package middleware

import (
	"errors"
	"net/http"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/Xlussov/EduCRM-be/pkg/response"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Caller = domain.Caller

func GetCaller(c echo.Context) (*Caller, error) {
	userToken := c.Get("user")
	if userToken == nil {
		return nil, errors.New("missing token")
	}

	token, ok := userToken.(*jwt.Token)
	if !ok {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	role := domain.Role(claims.Role)
	if !isKnownRole(role) {
		return nil, errors.New("invalid role")
	}

	branchIDs := make([]uuid.UUID, 0, len(claims.BranchIDs))
	for _, rawID := range claims.BranchIDs {
		id, err := uuid.Parse(rawID)
		if err != nil {
			return nil, errors.New("invalid branch id")
		}
		branchIDs = append(branchIDs, id)
	}

	return &Caller{
		UserID:    userID,
		Role:      role,
		BranchIDs: branchIDs,
	}, nil
}

func RequireRoles(roles ...string) echo.MiddlewareFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			caller, err := GetCaller(c)
			if err != nil {
				return response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error(), nil)
			}

			if caller.Role == domain.RoleSuperadmin {
				return next(c)
			}

			if len(allowed) > 0 {
				if _, ok := allowed[string(caller.Role)]; !ok {
					return response.Error(c, http.StatusForbidden, "ROLE_NOT_ALLOWED", "Role is not allowed", nil)
				}
			}

			return next(c)
		}
	}
}

func isKnownRole(role domain.Role) bool {
	switch role {
	case domain.RoleSuperadmin, domain.RoleAdmin, domain.RoleTeacher:
		return true
	default:
		return false
	}
}
