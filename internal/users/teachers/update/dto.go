package update

import "github.com/google/uuid"

type Request struct {
	FirstName string    `json:"first_name" validate:"required,min=2,max=50"`
	LastName  string    `json:"last_name" validate:"required,min=2,max=50"`
	Phone     string    `json:"phone" validate:"required,e164"`
	BranchID  uuid.UUID `json:"branch_id" validate:"required,uuid"`
}

type BranchResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type Response struct {
	ID        uuid.UUID        `json:"id"`
	FirstName string           `json:"first_name"`
	LastName  string           `json:"last_name"`
	Phone     string           `json:"phone"`
	Status    string           `json:"status"`
	Branches  []BranchResponse `json:"branches"`
}
