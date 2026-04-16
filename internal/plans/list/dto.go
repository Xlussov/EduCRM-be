package list

import "github.com/google/uuid"

type Request struct {
	BranchID uuid.UUID `query:"branch_id"`
}

type Subject struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type PricingGrid struct {
	Lessons int     `json:"lessons"`
	Price   float64 `json:"price"`
}

type PlanResponse struct {
	ID          uuid.UUID     `json:"id"`
	Name        string        `json:"name"`
	Type        string        `json:"type"`
	Subjects    []Subject     `json:"subjects"`
	PricingGrid []PricingGrid `json:"pricing_grid"`
	Status      string        `json:"status"`
}
