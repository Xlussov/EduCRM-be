package get

import (
	"time"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
)

type Request struct {
	ID uuid.UUID
}

type SubjectResponse struct {
	ID          string              `json:"id"`
	BranchID    string              `json:"branch_id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Status      domain.EntityStatus `json:"status"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

type Response struct {
	Subject SubjectResponse `json:"subject"`
}
