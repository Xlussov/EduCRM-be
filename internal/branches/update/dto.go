package update

import "time"

type Request struct {
	Name    string `json:"name" validate:"required,min=2,max=50"`
	Address string `json:"address" validate:"required,min=2,max=255"`
	City    string `json:"city" validate:"required,min=2,max=50"`
}

type Response struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	City      string    `json:"city"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
