package list

import "github.com/google/uuid"

type Request struct {
	BranchID *uuid.UUID `query:"branch_id"`
}

type BranchResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type TeacherResponse struct {
	ID        uuid.UUID        `json:"id"`
	FirstName string           `json:"first_name"`
	LastName  string           `json:"last_name"`
	Phone     string           `json:"phone"`
	Status    string           `json:"status"`
	Branches  []BranchResponse `json:"branches"`
}
