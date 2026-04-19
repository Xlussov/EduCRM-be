package list

import "github.com/google/uuid"

type Request struct {
	FromDate  string     `query:"from_date"`
	ToDate    string     `query:"to_date"`
	TeacherID *uuid.UUID `query:"teacher_id"`
	StudentID *uuid.UUID `query:"student_id"`
	GroupID   *uuid.UUID `query:"group_id"`
}

type TeacherRef struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
}

type SubjectRef struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type StudentRef struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
}

type GroupRef struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type LessonResponse struct {
	ID        uuid.UUID   `json:"id"`
	Date      string      `json:"date"`
	StartTime string      `json:"start_time"`
	EndTime   string      `json:"end_time"`
	Status    string      `json:"status"`
	Teacher   TeacherRef  `json:"teacher"`
	Subject   SubjectRef  `json:"subject"`
	Student   *StudentRef `json:"student,omitempty"`
	Group     *GroupRef   `json:"group,omitempty"`
}
