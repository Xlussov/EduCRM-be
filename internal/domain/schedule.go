package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type LessonStatus string

const (
	LessonStatusScheduled LessonStatus = "SCHEDULED"
	LessonStatusCompleted LessonStatus = "COMPLETED"
	LessonStatusCancelled LessonStatus = "CANCELLED"
)

type Lesson struct {
	ID         uuid.UUID
	BranchID   uuid.UUID
	TemplateID *uuid.UUID
	TeacherID  uuid.UUID
	SubjectID  uuid.UUID
	StudentID  *uuid.UUID
	GroupID    *uuid.UUID
	Date       time.Time
	StartTime  time.Time
	EndTime    time.Time
	Status     LessonStatus
	CreatedAt  time.Time
}

type LessonDetails struct {
	ID               uuid.UUID
	BranchID         uuid.UUID
	TemplateID       *uuid.UUID
	TeacherID        uuid.UUID
	TeacherFirstName string
	TeacherLastName  string
	SubjectID        uuid.UUID
	SubjectName      string
	StudentID        *uuid.UUID
	StudentFirstName string
	StudentLastName  string
	GroupID          *uuid.UUID
	GroupName        string
	Date             time.Time
	StartTime        time.Time
	EndTime          time.Time
	Status           LessonStatus
	CreatedAt        time.Time
}

type Template struct {
	ID         uuid.UUID
	BranchID   uuid.UUID
	TeacherID  uuid.UUID
	SubjectID  uuid.UUID
	StudentID  *uuid.UUID
	GroupID    *uuid.UUID
	DaysOfWeek []int32
	StartTime  time.Time
	EndTime    time.Time
	StartDate  time.Time
	EndDate    time.Time
	IsActive   bool
}

type ScheduleRepository interface {
	CreateLesson(ctx context.Context, lesson *Lesson) error
	CreateTemplate(ctx context.Context, template *Template) error
	BulkCreateLessons(ctx context.Context, lessons []Lesson) error
	UpdateLessonStatus(ctx context.Context, id uuid.UUID, status LessonStatus) error
	UpdateLesson(ctx context.Context, lesson *Lesson) error
	GetLessonByID(ctx context.Context, id uuid.UUID) (*Lesson, error)
	CheckTeacherConflict(ctx context.Context, teacherID uuid.UUID, date time.Time, start, end time.Time) (bool, error)
	CheckTeacherConflictExcludingLesson(ctx context.Context, teacherID uuid.UUID, date time.Time, start, end time.Time, lessonID uuid.UUID) (bool, error)
	CheckStudentConflict(ctx context.Context, studentID uuid.UUID, date time.Time, start, end time.Time) (bool, error)
	CheckStudentConflictExcludingLesson(ctx context.Context, studentID uuid.UUID, date time.Time, start, end time.Time, lessonID uuid.UUID) (bool, error)
	CheckTeacherFutureLessonsInBranch(ctx context.Context, teacherID, branchID uuid.UUID) (bool, error)
	CheckTeacherActiveTemplatesInBranch(ctx context.Context, teacherID, branchID uuid.UUID) (bool, error)
	ListLessons(ctx context.Context, from, to time.Time, teacherID, studentID, groupID *uuid.UUID, branchIDs []uuid.UUID) ([]LessonDetails, error)
	GetTeacherSchedule(ctx context.Context, teacherID uuid.UUID, from, to time.Time) ([]Lesson, error)
}
