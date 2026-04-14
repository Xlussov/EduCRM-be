package list

import (
	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
)

type Request struct {
	BranchID uuid.UUID `query:"branch_id"`
	Search   string    `query:"search"`
	Status   string    `query:"status"`
}

type StudentResponse struct {
	ID        uuid.UUID           `json:"id"`
	FirstName string              `json:"first_name"`
	LastName  string              `json:"last_name"`
	Phone     *string             `json:"phone"`
	Email     *string             `json:"email"`
	Status    domain.EntityStatus `json:"status"`
}

type Response struct {
	Students []StudentResponse `json:"students"`
}
