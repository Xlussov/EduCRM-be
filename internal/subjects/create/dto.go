package create

import "github.com/google/uuid"

type Request struct {
	BranchID    uuid.UUID `json:"branch_id" validate:"required"`
	Name        string    `json:"name" validate:"required,min=2,max=50"`
	Description string    `json:"description" validate:"omitempty,min=2,max=500"`
}

type Response struct {
	ID          string `json:"id"`
	BranchID    string `json:"branch_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
