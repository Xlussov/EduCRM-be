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

type UserBranch struct {
	ID   uuid.UUID
	Name string
}

type UserWithBranches struct {
	ID        uuid.UUID
	Phone     string
	FirstName string
	LastName  string
	Role      Role
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
	Branches  []UserBranch
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByPhone(ctx context.Context, phone string) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetWithBranchesByID(ctx context.Context, id uuid.UUID) (*UserWithBranches, error)
	GetAdmins(ctx context.Context) ([]*UserWithBranches, error)
	GetTeachers(ctx context.Context, branchIDs []uuid.UUID) ([]*UserWithBranches, error)
	UpdateUser(ctx context.Context, user *User) error
	UpdateUserStatus(ctx context.Context, id uuid.UUID, isActive bool) error
	DeleteUserBranches(ctx context.Context, userID uuid.UUID) error
	AssignToBranches(ctx context.Context, userID uuid.UUID, branchIDs []uuid.UUID) error
	GetUserBranchIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
	IsBranchActive(ctx context.Context, branchID uuid.UUID) (bool, error)
	CountActiveBranchesByIDs(ctx context.Context, branchIDs []uuid.UUID) (int, error)
}
