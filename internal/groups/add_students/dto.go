package addstudents

import "github.com/google/uuid"

type Request struct {
	StudentIDs []uuid.UUID `json:"student_ids" validate:"required,min=1"`
}

type Response struct {
	Message string `json:"message"`
}
