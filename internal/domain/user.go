package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleSuperadmin Role = "SUPERADMIN"
	RoleAdmin      Role = "ADMIN"
	RoleTeacher    Role = "TEACHER"
)

type User struct {
	ID           uuid.UUID
	Phone        string
	PasswordHash string
	FirstName    string
	LastName     string
	Role         Role
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByPhone(ctx context.Context, phone string) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	AssignToBranches(ctx context.Context, userID uuid.UUID, branchIDs []uuid.UUID) error
	GetUserBranchIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
}
