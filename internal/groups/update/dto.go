package update

import (
	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
)

type Request struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
}

type Response struct {
	ID       uuid.UUID           `json:"id"`
	BranchID uuid.UUID           `json:"branch_id"`
	Name     string              `json:"name"`
	Status   domain.EntityStatus `json:"status"`
}
