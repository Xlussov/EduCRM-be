package teachers

import "github.com/google/uuid"

type Request struct {
	Phone     string    `json:"phone"`
	Password  string    `json:"password"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	BranchID  uuid.UUID `json:"branch_id"`
}

type Response struct {
	ID uuid.UUID `json:"id"`
}
