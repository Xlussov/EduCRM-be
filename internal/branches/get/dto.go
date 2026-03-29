package get

import "github.com/google/uuid"

type Response struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Address string    `json:"address"`
	City    string    `json:"city"`
	Status  string    `json:"status"`
}
