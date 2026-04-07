package list

import (
	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
)

type Request struct {
	BranchID uuid.UUID `query:"branch_id"`
	Status   string    `query:"status"`
}

type GroupResponse struct {
	ID            uuid.UUID           `json:"id"`
	Name          string              `json:"name"`
	StudentsCount int                 `json:"students_count"`
	Status        domain.EntityStatus `json:"status"`
}

type Response struct {
	Groups []GroupResponse `json:"groups"`
}
