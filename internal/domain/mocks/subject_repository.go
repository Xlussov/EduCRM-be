package mocks

import (
	"context"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/stretchr/testify/mock"
)

type SubjectRepository struct {
	mock.Mock
}

func (m *SubjectRepository) Create(ctx context.Context, subject *domain.Subject) error {
	args := m.Called(ctx, subject)
	return args.Error(0)
}
