package create

import (
	"time"

	"github.com/google/uuid"
)

type Request struct {
	PlanID    uuid.UUID `json:"plan_id" validate:"required"`
	SubjectID uuid.UUID `json:"subject_id" validate:"required"`
	StartDate time.Time `json:"start_date" validate:"required"`
}

type Response struct {
	ID uuid.UUID `json:"id"`
}
