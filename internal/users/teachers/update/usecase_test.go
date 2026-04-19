package update

import (
	"context"
	"testing"
	"time"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/Xlussov/EduCRM-be/internal/domain/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUseCase_Execute(t *testing.T) {
	teacherID := uuid.New()
	currentBranchID := uuid.New()
	newBranchID := uuid.New()
	otherBranchID := uuid.New()

	req := Request{
		FirstName: "Updated",
		LastName:  "Teacher",
		Phone:     "+998901234567",
		BranchIDs: []uuid.UUID{currentBranchID, newBranchID},
	}

	tests := []struct {
		name          string
		caller        domain.Caller
		mockSetup     func(repo *mocks.UserRepository, scheduleRepo *mocks.ScheduleRepository)
		expectedError error
	}{
		{
			name:   "superadmin_success",
			caller: domain.Caller{Role: domain.RoleSuperadmin},
			mockSetup: func(repo *mocks.UserRepository, scheduleRepo *mocks.ScheduleRepository) {
				repo.On("GetWithBranchesByID", mock.Anything, teacherID).Return(&domain.UserWithBranches{
					ID:       teacherID,
					Role:     domain.RoleTeacher,
					IsActive: true,
					Branches: []domain.UserBranch{{ID: currentBranchID, Name: "Current"}},
				}, nil).Once()
				repo.On("CountActiveBranchesByIDs", mock.Anything, []uuid.UUID{currentBranchID, newBranchID}).Return(2, nil).Once()
				repo.On("UpdateUser", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
					return u.ID == teacherID && u.Phone == req.Phone && u.FirstName == req.FirstName && u.LastName == req.LastName
				})).Return(nil).Once()
				repo.On("DeleteUserBranches", mock.Anything, teacherID).Return(nil).Once()
				repo.On("AssignToBranches", mock.Anything, teacherID, []uuid.UUID{currentBranchID, newBranchID}).Return(nil).Once()
				repo.On("GetWithBranchesByID", mock.Anything, teacherID).Return(&domain.UserWithBranches{
					ID:        teacherID,
					Phone:     req.Phone,
					FirstName: req.FirstName,
					LastName:  req.LastName,
					Role:      domain.RoleTeacher,
					IsActive:  true,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Branches:  []domain.UserBranch{{ID: currentBranchID, Name: "Current"}, {ID: newBranchID, Name: "New"}},
				}, nil).Once()
			},
			expectedError: nil,
		},
		{
			name:   "admin_no_current_access",
			caller: domain.Caller{Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{otherBranchID}},
			mockSetup: func(repo *mocks.UserRepository, scheduleRepo *mocks.ScheduleRepository) {
				repo.On("GetWithBranchesByID", mock.Anything, teacherID).Return(&domain.UserWithBranches{
					ID:       teacherID,
					Role:     domain.RoleTeacher,
					IsActive: true,
					Branches: []domain.UserBranch{{ID: currentBranchID, Name: "Current"}},
				}, nil).Once()
			},
			expectedError: domain.ErrBranchAccessDenied,
		},
		{
			name:   "admin_no_new_access",
			caller: domain.Caller{Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{currentBranchID}},
			mockSetup: func(repo *mocks.UserRepository, scheduleRepo *mocks.ScheduleRepository) {
				repo.On("GetWithBranchesByID", mock.Anything, teacherID).Return(&domain.UserWithBranches{
					ID:       teacherID,
					Role:     domain.RoleTeacher,
					IsActive: true,
					Branches: []domain.UserBranch{{ID: currentBranchID, Name: "Current"}},
				}, nil).Once()
			},
			expectedError: domain.ErrBranchAccessDenied,
		},
		{
			name:   "archived_teacher",
			caller: domain.Caller{Role: domain.RoleSuperadmin},
			mockSetup: func(repo *mocks.UserRepository, scheduleRepo *mocks.ScheduleRepository) {
				repo.On("GetWithBranchesByID", mock.Anything, teacherID).Return(&domain.UserWithBranches{
					ID:       teacherID,
					Role:     domain.RoleTeacher,
					IsActive: false,
				}, nil).Once()
			},
			expectedError: domain.ErrCannotEditArchived,
		},
		{
			name:   "not_teacher",
			caller: domain.Caller{Role: domain.RoleSuperadmin},
			mockSetup: func(repo *mocks.UserRepository, scheduleRepo *mocks.ScheduleRepository) {
				repo.On("GetWithBranchesByID", mock.Anything, teacherID).Return(&domain.UserWithBranches{
					ID:   teacherID,
					Role: domain.RoleAdmin,
				}, nil).Once()
			},
			expectedError: domain.ErrNotFound,
		},
		{
			name:   "archived_target_branch",
			caller: domain.Caller{Role: domain.RoleSuperadmin},
			mockSetup: func(repo *mocks.UserRepository, scheduleRepo *mocks.ScheduleRepository) {
				repo.On("GetWithBranchesByID", mock.Anything, teacherID).Return(&domain.UserWithBranches{
					ID:       teacherID,
					Role:     domain.RoleTeacher,
					IsActive: true,
					Branches: []domain.UserBranch{{ID: currentBranchID, Name: "Current"}},
				}, nil).Once()
				repo.On("CountActiveBranchesByIDs", mock.Anything, []uuid.UUID{currentBranchID, newBranchID}).Return(1, nil).Once()
			},
			expectedError: domain.ErrArchivedReference,
		},
		{
			name:   "phone_conflict",
			caller: domain.Caller{Role: domain.RoleSuperadmin},
			mockSetup: func(repo *mocks.UserRepository, scheduleRepo *mocks.ScheduleRepository) {
				repo.On("GetWithBranchesByID", mock.Anything, teacherID).Return(&domain.UserWithBranches{
					ID:       teacherID,
					Role:     domain.RoleTeacher,
					IsActive: true,
					Branches: []domain.UserBranch{{ID: currentBranchID, Name: "Current"}},
				}, nil).Once()
				repo.On("CountActiveBranchesByIDs", mock.Anything, []uuid.UUID{currentBranchID, newBranchID}).Return(2, nil).Once()
				repo.On("UpdateUser", mock.Anything, mock.Anything).Return(domain.ErrAlreadyExists).Once()
			},
			expectedError: domain.ErrPhoneAlreadyExists,
		},
		{
			name:   "removed_branch_has_future_lessons",
			caller: domain.Caller{Role: domain.RoleSuperadmin},
			mockSetup: func(repo *mocks.UserRepository, scheduleRepo *mocks.ScheduleRepository) {
				repo.On("GetWithBranchesByID", mock.Anything, teacherID).Return(&domain.UserWithBranches{
					ID:       teacherID,
					Role:     domain.RoleTeacher,
					IsActive: true,
					Branches: []domain.UserBranch{{ID: currentBranchID, Name: "Current"}, {ID: otherBranchID, Name: "Other"}},
				}, nil).Once()
				scheduleRepo.On("CheckTeacherFutureLessonsInBranch", mock.Anything, teacherID, otherBranchID).Return(true, nil).Once()
			},
			expectedError: domain.ErrTeacherHasFutureLessons,
		},
		{
			name:   "removed_branch_has_active_templates",
			caller: domain.Caller{Role: domain.RoleSuperadmin},
			mockSetup: func(repo *mocks.UserRepository, scheduleRepo *mocks.ScheduleRepository) {
				repo.On("GetWithBranchesByID", mock.Anything, teacherID).Return(&domain.UserWithBranches{
					ID:       teacherID,
					Role:     domain.RoleTeacher,
					IsActive: true,
					Branches: []domain.UserBranch{{ID: currentBranchID, Name: "Current"}, {ID: otherBranchID, Name: "Other"}},
				}, nil).Once()
				scheduleRepo.On("CheckTeacherFutureLessonsInBranch", mock.Anything, teacherID, otherBranchID).Return(false, nil).Once()
				scheduleRepo.On("CheckTeacherActiveTemplatesInBranch", mock.Anything, teacherID, otherBranchID).Return(true, nil).Once()
			},
			expectedError: domain.ErrTeacherHasActiveTemplates,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.UserRepository)
			scheduleRepo := new(mocks.ScheduleRepository)
			tt.mockSetup(repo, scheduleRepo)

			uc := NewUseCase(repo, scheduleRepo, &mocks.MockTxManager{})
			res, err := uc.Execute(context.Background(), tt.caller, teacherID, req)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, teacherID, res.ID)
			}

			repo.AssertExpectations(t)
			scheduleRepo.AssertExpectations(t)
		})
	}
}
