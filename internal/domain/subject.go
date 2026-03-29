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
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type SubjectRepository interface {
	Create(ctx context.Context, subject *Subject) error
}
