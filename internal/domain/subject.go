package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Subject struct {
	ID          uuid.UUID
	Name        string
	Description string
	Status      EntityStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type SubjectRepository interface {
	Create(ctx context.Context, subject *Subject) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status EntityStatus) error
	GetAll(ctx context.Context) ([]*Subject, error)
	Update(ctx context.Context, subject *Subject) error
}
