package create

import "github.com/google/uuid"

type Request struct {
	Phone     string      `json:"phone" validate:"required,e164"`
	Password  string      `json:"password" validate:"required,min=6"`
	FirstName string      `json:"first_name" validate:"required,min=2,max=50"`
	LastName  string      `json:"last_name" validate:"required,min=2,max=50"`
	BranchIDs []uuid.UUID `json:"branch_ids" validate:"required,min=1,unique,dive,uuid"`
}

type Response struct {
	ID uuid.UUID `json:"id"`
}
