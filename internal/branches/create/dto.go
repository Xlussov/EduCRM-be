package create

import "github.com/google/uuid"

type Request struct {
	Name    string `json:"name" validate:"required"`
	Address string `json:"address" validate:"required"`
	City    string `json:"city" validate:"required"`
}

type Response struct {
	ID uuid.UUID `json:"id"`
}
