package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Student struct {
	ID                 uuid.UUID
	BranchID           uuid.UUID
	FirstName          string
	LastName           string
	Dob                *time.Time
	Phone              *string
	Email              *string
	Address            *string
	ParentName         string
	ParentPhone        string
	ParentEmail        *string
	ParentRelationship *string
	Status             EntityStatus
	CreatedAt          time.Time
}

type StudentRepository interface {
	Create(ctx context.Context, student *Student) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status EntityStatus) error
	GetBranchID(ctx context.Context, id uuid.UUID) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Student, error)
	Update(ctx context.Context, student *Student) (*Student, error)
	GetByBranchID(ctx context.Context, branchID uuid.UUID) ([]*Student, error)
}
