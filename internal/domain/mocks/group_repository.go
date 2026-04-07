package mocks

import (
	"context"
	"time"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type GroupRepository struct {
	mock.Mock
}

func (m *GroupRepository) Create(ctx context.Context, group *domain.Group) error {
	args := m.Called(ctx, group)
	return args.Error(0)
}

func (m *GroupRepository) GetByBranchID(ctx context.Context, branchID uuid.UUID, status *domain.EntityStatus) ([]*domain.GroupWithCount, error) {
	args := m.Called(ctx, branchID, status)
	var r0 []*domain.GroupWithCount
	if args.Get(0) != nil {
		r0 = args.Get(0).([]*domain.GroupWithCount)
	}
	return r0, args.Error(1)
}

func (m *GroupRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Group, error) {
	args := m.Called(ctx, id)
	var r0 *domain.Group
	if args.Get(0) != nil {
		r0 = args.Get(0).(*domain.Group)
	}
	return r0, args.Error(1)
}

func (m *GroupRepository) UpdateName(ctx context.Context, id uuid.UUID, name string) (*domain.Group, error) {
	args := m.Called(ctx, id, name)
	var r0 *domain.Group
	if args.Get(0) != nil {
		r0 = args.Get(0).(*domain.Group)
	}
	return r0, args.Error(1)
}

func (m *GroupRepository) AddStudent(ctx context.Context, groupID, studentID uuid.UUID, joinedAt time.Time) error {
	args := m.Called(ctx, groupID, studentID, joinedAt)
	return args.Error(0)
}

func (m *GroupRepository) RemoveStudent(ctx context.Context, groupID, studentID uuid.UUID, leftAt time.Time) error {
	args := m.Called(ctx, groupID, studentID, leftAt)
	return args.Error(0)
}

func (m *GroupRepository) GetActiveStudentIDs(ctx context.Context, groupID uuid.UUID) ([]uuid.UUID, error) {
	args := m.Called(ctx, groupID)
	var r0 []uuid.UUID
	if args.Get(0) != nil {
		r0 = args.Get(0).([]uuid.UUID)
	}
	return r0, args.Error(1)
}

func (m *GroupRepository) GetStudents(ctx context.Context, groupID uuid.UUID) ([]*domain.GroupStudent, error) {
	args := m.Called(ctx, groupID)
	var r0 []*domain.GroupStudent
	if args.Get(0) != nil {
		r0 = args.Get(0).([]*domain.GroupStudent)
	}
	return r0, args.Error(1)
}

func (m *GroupRepository) GetBranchID(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	args := m.Called(ctx, id)
	var r0 uuid.UUID
	if args.Get(0) != nil {
		r0 = args.Get(0).(uuid.UUID)
	}
	return r0, args.Error(1)
}
