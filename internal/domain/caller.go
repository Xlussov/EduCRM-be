package domain

import "github.com/google/uuid"

type Caller struct {
	UserID    uuid.UUID
	Role      Role
	BranchIDs []uuid.UUID
}
