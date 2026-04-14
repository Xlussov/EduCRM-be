package syncstudents

import "github.com/google/uuid"

type Request struct {
	StudentIDs []uuid.UUID `json:"student_ids"`
}

type Response struct {
	Message string `json:"message"`
}
