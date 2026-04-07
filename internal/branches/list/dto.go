package list

import "github.com/google/uuid"

type Request struct {
	Status string `query:"status"`
}

type BranchResponse struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Address string    `json:"address"`
	City    string    `json:"city"`
	Status  string    `json:"status"`
}
