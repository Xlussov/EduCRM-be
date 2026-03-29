package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type EntityStatus string

const (
	StatusActive   EntityStatus = "ACTIVE"
	StatusArchived EntityStatus = "ARCHIVED"
)

type Branch struct {
	ID        uuid.UUID
	Name      string
	Address   string
	City      string
	Status    EntityStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

type BranchRepository interface {
	Create(ctx context.Context, branch *Branch) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status EntityStatus) error
	GetAll(ctx context.Context) ([]*Branch, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*Branch, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Branch, error)
	Update(ctx context.Context, branch *Branch) error
}
