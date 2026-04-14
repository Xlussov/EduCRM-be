package get

import (
	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
)

type StudentResponse struct {
	ID        uuid.UUID           `json:"id"`
	FirstName string              `json:"first_name"`
	LastName  string              `json:"last_name"`
	Status    domain.EntityStatus `json:"status"`
	Phone     *string             `json:"phone"`
	Email     *string             `json:"email"`
}

type Response struct {
	ID       uuid.UUID           `json:"id"`
	Name     string              `json:"name"`
	Status   domain.EntityStatus `json:"status"`
	BranchID uuid.UUID           `json:"branch_id"`
	Students []StudentResponse   `json:"students"`
}
