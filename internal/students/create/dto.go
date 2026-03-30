package create

import "github.com/google/uuid"

type Request struct {
	BranchID           uuid.UUID `json:"branch_id" validate:"required"`
	FirstName          string    `json:"first_name" validate:"required"`
	LastName           string    `json:"last_name" validate:"required"`
	Dob                *string   `json:"dob"`
	Phone              *string   `json:"phone"`
	Email              *string   `json:"email"`
	Address            *string   `json:"address"`
	ParentName         string    `json:"parent_name" validate:"required"`
	ParentPhone        string    `json:"parent_phone" validate:"required"`
	ParentEmail        *string   `json:"parent_email"`
	ParentRelationship *string   `json:"parent_relationship"`
}

type Response struct {
	ID uuid.UUID `json:"id"`
}
