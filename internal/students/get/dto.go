package get

import (
	"time"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
)

type Response struct {
	ID                 uuid.UUID           `json:"id"`
	BranchID           uuid.UUID           `json:"branch_id"`
	FirstName          string              `json:"first_name"`
	LastName           string              `json:"last_name"`
	Dob                *string             `json:"dob,omitempty"`
	Phone              *string             `json:"phone,omitempty"`
	Email              *string             `json:"email,omitempty"`
	Address            *string             `json:"address,omitempty"`
	ParentName         string              `json:"parent_name"`
	ParentPhone        string              `json:"parent_phone"`
	ParentEmail        *string             `json:"parent_email,omitempty"`
	ParentRelationship *string             `json:"parent_relationship,omitempty"`
	Status             domain.EntityStatus `json:"status"`
	CreatedAt          time.Time           `json:"created_at"`
}
