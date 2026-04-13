package get

import "github.com/google/uuid"

type Request struct {
	ID uuid.UUID
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
