package create

import "github.com/google/uuid"

type Request struct {
	BranchID uuid.UUID `json:"branch_id" validate:"required"`
	Name     string    `json:"name" validate:"required,min=1,max=255"`
}

type Response struct {
	ID uuid.UUID `json:"id"`
}
