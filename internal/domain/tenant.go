package domain

import "github.com/google/uuid"

func RequiresBranchAccess(role Role) bool {
	switch role {
	case RoleAdmin, RoleTeacher:
		return true
	default:
		return false
	}
}

func HasBranchAccess(branchIDs []uuid.UUID, branchID uuid.UUID) bool {
	for _, id := range branchIDs {
		if id == branchID {
			return true
		}
	}
	return false
}
