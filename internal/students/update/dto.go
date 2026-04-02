package update

import "time"

type Request struct {
	FirstName          string  `json:"first_name" validate:"required,min=2,max=50"`
	LastName           string  `json:"last_name" validate:"required,min=2,max=50"`
	Dob                *string `json:"dob" validate:"omitempty,datetime=2006-01-02"`
	Phone              *string `json:"phone" validate:"omitempty,e164"`
	Email              *string `json:"email" validate:"omitempty,email"`
	Address            *string `json:"address" validate:"omitempty,min=2,max=255"`
	ParentName         string  `json:"parent_name" validate:"required,min=2,max=50"`
	ParentPhone        string  `json:"parent_phone" validate:"required,e164"`
	ParentEmail        *string `json:"parent_email" validate:"omitempty,email"`
	ParentRelationship *string `json:"parent_relationship" validate:"omitempty,min=2,max=50"`
}

type Response struct {
	ID                 string    `json:"id"`
	BranchID           string    `json:"branch_id"`
	FirstName          string    `json:"first_name"`
	LastName           string    `json:"last_name"`
	Dob                *string   `json:"dob"`
	Phone              *string   `json:"phone"`
	Email              *string   `json:"email"`
	Address            *string   `json:"address"`
	ParentName         string    `json:"parent_name"`
	ParentPhone        string    `json:"parent_phone"`
	ParentEmail        *string   `json:"parent_email"`
	ParentRelationship *string   `json:"parent_relationship"`
	Status             string    `json:"status"`
	CreatedAt          time.Time `json:"created_at"`
}
