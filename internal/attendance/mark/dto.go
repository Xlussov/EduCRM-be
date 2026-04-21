package mark

import "github.com/google/uuid"

type AttendanceItem struct {
	StudentID uuid.UUID `json:"student_id"`
	IsPresent *bool     `json:"is_present"`
	Notes     *string   `json:"notes"`
}

type Request struct {
	Attendance []AttendanceItem `json:"attendance"`
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
