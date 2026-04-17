package cancel

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/Xlussov/EduCRM-be/internal/domain/mocks"
)

func TestUseCase_Execute(t *testing.T) {
	userID := uuid.New()
	targetBranchID := uuid.New()
	lessonID := uuid.New()

	lesson := &domain.Lesson{
		ID:       lessonID,
		BranchID: targetBranchID,
	}

	tests := []struct {
		name        string
		lessonID    uuid.UUID
		caller      domain.Caller
		mockSetup   func(sRepo *mocks.ScheduleRepository)
		expectedErr error
	}{
		{
			name:     "success (superadmin)",
			lessonID: lessonID,
			caller:   domain.Caller{UserID: userID, Role: domain.RoleSuperadmin, BranchIDs: []uuid.UUID{}},
			mockSetup: func(sRepo *mocks.ScheduleRepository) {
				sRepo.On("GetLessonByID", mock.Anything, lessonID).Return(lesson, nil).Once()
				sRepo.On("UpdateLessonStatus", mock.Anything, lessonID, domain.LessonStatusCancelled).Return(nil).Once()
			},
			expectedErr: nil,
		},
		{
			name:     "success (admin with branch access)",
			lessonID: lessonID,
			caller:   domain.Caller{UserID: userID, Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{targetBranchID}},
			mockSetup: func(sRepo *mocks.ScheduleRepository) {
				sRepo.On("GetLessonByID", mock.Anything, lessonID).Return(lesson, nil).Once()
				sRepo.On("UpdateLessonStatus", mock.Anything, lessonID, domain.LessonStatusCancelled).Return(nil).Once()
			},
			expectedErr: nil,
		},
		{
			name:     "admin branch access denied",
			lessonID: lessonID,
			caller:   domain.Caller{UserID: userID, Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{uuid.New()}},
			mockSetup: func(sRepo *mocks.ScheduleRepository) {
				sRepo.On("GetLessonByID", mock.Anything, lessonID).Return(lesson, nil).Once()
			},
			expectedErr: domain.ErrBranchAccessDenied,
		},
		{
			name:     "lesson not found",
			lessonID: lessonID,
			caller:   domain.Caller{UserID: userID, Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{targetBranchID}},
			mockSetup: func(sRepo *mocks.ScheduleRepository) {
				sRepo.On("GetLessonByID", mock.Anything, lessonID).Return(nil, errors.New("not found")).Once()
			},
			expectedErr: errors.New("not found"),
		},
		{
			name:     "error updating status",
			lessonID: lessonID,
			caller:   domain.Caller{UserID: userID, Role: domain.RoleSuperadmin, BranchIDs: []uuid.UUID{}},
			mockSetup: func(sRepo *mocks.ScheduleRepository) {
				sRepo.On("GetLessonByID", mock.Anything, lessonID).Return(lesson, nil).Once()
				sRepo.On("UpdateLessonStatus", mock.Anything, lessonID, domain.LessonStatusCancelled).Return(errors.New("db error")).Once()
			},
			expectedErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sRepo := new(mocks.ScheduleRepository)
			if tt.mockSetup != nil {
				tt.mockSetup(sRepo)
			}

			uc := NewUseCase(sRepo)
			res, err := uc.Execute(context.Background(), tt.caller, tt.lessonID)

			if tt.expectedErr != nil {
				if tt.expectedErr.Error() == "not found" || tt.expectedErr.Error() == "db error" {
					assert.Error(t, err)
				} else {
					assert.ErrorIs(t, err, tt.expectedErr)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "success", res.Message)
			}

			sRepo.AssertExpectations(t)
		})
	}
}
