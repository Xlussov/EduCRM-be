package list

import (
	"time"

	"github.com/google/uuid"
)

type SubscriptionResponse struct {
	ID        uuid.UUID  `json:"id"`
	Plan      PlanRef    `json:"plan"`
	Subject   SubjectRef `json:"subject"`
	StartDate time.Time  `json:"start_date"`
	CreatedAt time.Time  `json:"created_at"`
}

type PlanRef struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type SubjectRef struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
