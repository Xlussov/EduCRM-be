package mocks

import (
	"context"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByPhone(ctx context.Context, phone string) (*domain.User, error) {
	args := m.Called(ctx, phone)
	var user *domain.User
	if args.Get(0) != nil {
		user = args.Get(0).(*domain.User)
	}
	return user, args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	args := m.Called(ctx, id)
	var user *domain.User
	if args.Get(0) != nil {
		user = args.Get(0).(*domain.User)
	}
	return user, args.Error(1)
}

func (m *MockUserRepository) GetWithBranchesByID(ctx context.Context, id uuid.UUID) (*domain.UserWithBranches, error) {
	args := m.Called(ctx, id)
	var user *domain.UserWithBranches
	if args.Get(0) != nil {
		user = args.Get(0).(*domain.UserWithBranches)
	}
	return user, args.Error(1)
}

func (m *MockUserRepository) GetAdmins(ctx context.Context) ([]*domain.UserWithBranches, error) {
	args := m.Called(ctx)
	var admins []*domain.UserWithBranches
	if args.Get(0) != nil {
		admins = args.Get(0).([]*domain.UserWithBranches)
	}
	return admins, args.Error(1)
}

func (m *MockUserRepository) GetTeachers(ctx context.Context, branchIDs []uuid.UUID) ([]*domain.UserWithBranches, error) {
	args := m.Called(ctx, branchIDs)
	var teachers []*domain.UserWithBranches
	if args.Get(0) != nil {
		teachers = args.Get(0).([]*domain.UserWithBranches)
	}
	return teachers, args.Error(1)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateUserStatus(ctx context.Context, id uuid.UUID, isActive bool) error {
	args := m.Called(ctx, id, isActive)
	return args.Error(0)
}

func (m *MockUserRepository) DeleteUserBranches(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUserRepository) AssignToBranches(ctx context.Context, userID uuid.UUID, branchIDs []uuid.UUID) error {
	args := m.Called(ctx, userID, branchIDs)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserBranchIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	args := m.Called(ctx, userID)
	var branchIDs []uuid.UUID
	if args.Get(0) != nil {
		branchIDs = args.Get(0).([]uuid.UUID)
	}
	return branchIDs, args.Error(1)
}

func (m *MockUserRepository) IsBranchActive(ctx context.Context, branchID uuid.UUID) (bool, error) {
	args := m.Called(ctx, branchID)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) CountActiveBranchesByIDs(ctx context.Context, branchIDs []uuid.UUID) (int, error) {
	args := m.Called(ctx, branchIDs)
	return args.Int(0), args.Error(1)
}

func (m *MockUserRepository) CheckTeacherInBranch(ctx context.Context, teacherID, branchID uuid.UUID) (bool, error) {
	args := m.Called(ctx, teacherID, branchID)
	return args.Bool(0), args.Error(1)
}
