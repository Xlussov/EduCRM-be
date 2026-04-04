package create

import "github.com/google/uuid"

type PricingGridItem struct {
	Lessons int     `json:"lessons" validate:"required,min=1"`
	Price   float64 `json:"price" validate:"required,min=0"`
}

type Request struct {
	BranchID    uuid.UUID         `json:"branch_id" validate:"required"`
	Name        string            `json:"name" validate:"required"`
	Type        string            `json:"type" validate:"required,oneof=INDIVIDUAL GROUP"`
	SubjectIDs  []uuid.UUID       `json:"subject_ids" validate:"required,min=1"`
	PricingGrid []PricingGridItem `json:"pricing_grid" validate:"required,min=1,dive"`
}

type Response struct {
	ID uuid.UUID `json:"id"`
}
