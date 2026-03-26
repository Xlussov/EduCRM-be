package admins

import "github.com/google/uuid"

type Request struct {
	Phone     string      `json:"phone" validate:"required"`
	Password  string      `json:"password" validate:"required"`
	FirstName string      `json:"first_name" validate:"required"`
	LastName  string      `json:"last_name" validate:"required"`
	BranchIDs []uuid.UUID `json:"branch_ids" validate:"required"`
}

type Response struct {
	ID uuid.UUID `json:"id"`
}
