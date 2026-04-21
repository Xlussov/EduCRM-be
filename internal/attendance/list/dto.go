package list

import "github.com/google/uuid"

type Request struct {
	ID uuid.UUID
}

type StudentAttendance struct {
	StudentID uuid.UUID `json:"student_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Status    string    `json:"status"`
	IsPresent *bool     `json:"is_present"`
	Notes     *string   `json:"notes"`
}

type Response struct {
	Attendance []StudentAttendance `json:"attendance"`
}
