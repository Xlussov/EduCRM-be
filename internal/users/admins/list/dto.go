package list

import "github.com/google/uuid"

type Request struct{}

type BranchResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type AdminResponse struct {
	ID        uuid.UUID        `json:"id"`
	FirstName string           `json:"first_name"`
	LastName  string           `json:"last_name"`
	Phone     string           `json:"phone"`
	Status    string           `json:"status"`
	Branches  []BranchResponse `json:"branches"`
}
