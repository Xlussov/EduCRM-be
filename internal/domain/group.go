package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Group struct {
	ID        uuid.UUID
	BranchID  uuid.UUID
	Name      string
	Status    EntityStatus
	CreatedAt time.Time
}

type GroupWithCount struct {
	Group
	StudentsCount int
}

type GroupStudent struct {
	ID        uuid.UUID
	FirstName string
	LastName  string
	Status    EntityStatus
	Phone     *string
	Email     *string
}

type GroupRepository interface {
	Create(ctx context.Context, group *Group) error
	GetByBranchID(ctx context.Context, branchID uuid.UUID, status *EntityStatus) ([]*GroupWithCount, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Group, error)
	UpdateName(ctx context.Context, id uuid.UUID, name string) (*Group, error)
	AddStudents(ctx context.Context, groupID uuid.UUID, studentIDs []uuid.UUID, joinedAt time.Time) error
	RemoveStudents(ctx context.Context, groupID uuid.UUID, studentIDs []uuid.UUID, leftAt time.Time) error
	GetActiveStudentIDs(ctx context.Context, groupID uuid.UUID) ([]uuid.UUID, error)
	GetStudents(ctx context.Context, groupID uuid.UUID) ([]*GroupStudent, error)
	GetBranchID(ctx context.Context, id uuid.UUID) (uuid.UUID, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status EntityStatus) error
}
