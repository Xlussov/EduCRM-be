package get

import (
	"github.com/google/uuid"
)

type StudentResponse struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
}

type Response struct {
	ID       uuid.UUID         `json:"id"`
	Name     string            `json:"name"`
	Students []StudentResponse `json:"students"`
}
