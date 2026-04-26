package teacher_statistics

import (
	"context"
	"testing"
	"time"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/Xlussov/EduCRM-be/internal/domain/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUseCase_Execute(t *testing.T) {
	teacherID := uuid.New()
	callerTeacherID := uuid.New()
	callerAdminID := uuid.New()
	sharedBranchID := uuid.New()
	now := time.Now()

	tests := []struct {
		name           string
		req            Request
		caller         domain.Caller
		mockBehavior   func(mockReportRepo *mocks.MockReportRepository, mockUserRepo *mocks.MockUserRepository)
		expectedResult *Response
		expectedError  error
	}{
		{
			name: "Success - Teacher role overrides req.TeacherID",
			req: Request{
				TeacherID: uuid.New(), // some other ID
				StartDate: &now,
				EndDate:   &now,
			},
			caller: domain.Caller{
				UserID:    callerTeacherID,
				Role:      domain.RoleTeacher,
				BranchIDs: []uuid.UUID{sharedBranchID},
			},
			mockBehavior: func(mockReportRepo *mocks.MockReportRepository, mockUserRepo *mocks.MockUserRepository) {
				mockReportRepo.On("GetTeacherStatistics", mock.Anything, callerTeacherID, &now, &now).
					Return(&domain.TeacherStatisticsReport{
						ScheduledLessons: 5,
						CompletedLessons: 2,
						CancelledLessons: 1,
					}, nil)
			},
			expectedResult: &Response{
				ScheduledLessons: 5,
				CompletedLessons: 2,
				CancelledLessons: 1,
			},
			expectedError: nil,
		},
		{
			name: "Success - Admin with shared branch access",
			req: Request{
				TeacherID: teacherID,
				StartDate: &now,
				EndDate:   &now,
			},
			caller: domain.Caller{
				UserID:    callerAdminID,
				Role:      domain.RoleAdmin,
				BranchIDs: []uuid.UUID{sharedBranchID, uuid.New()},
			},
			mockBehavior: func(mockReportRepo *mocks.MockReportRepository, mockUserRepo *mocks.MockUserRepository) {
				mockUserRepo.On("GetUserBranchIDs", mock.Anything, teacherID).
					Return([]uuid.UUID{sharedBranchID}, nil)

				mockReportRepo.On("GetTeacherStatistics", mock.Anything, teacherID, &now, &now).
					Return(&domain.TeacherStatisticsReport{
						ScheduledLessons: 10,
						CompletedLessons: 8,
						CancelledLessons: 0,
					}, nil)
			},
			expectedResult: &Response{
				ScheduledLessons: 10,
				CompletedLessons: 8,
				CancelledLessons: 0,
			},
			expectedError: nil,
		},
		{
			name: "Error - Admin cross-branch access denied",
			req: Request{
				TeacherID: teacherID,
			},
			caller: domain.Caller{
				UserID:    callerAdminID,
				Role:      domain.RoleAdmin,
				BranchIDs: []uuid.UUID{uuid.New()}, // Different branch
			},
			mockBehavior: func(mockReportRepo *mocks.MockReportRepository, mockUserRepo *mocks.MockUserRepository) {
				mockUserRepo.On("GetUserBranchIDs", mock.Anything, teacherID).
					Return([]uuid.UUID{sharedBranchID}, nil)
			},
			expectedResult: nil,
			expectedError:  domain.ErrBranchAccessDenied,
		},
		{
			name: "Error - Admin requesting teacher with no branches",
			req: Request{
				TeacherID: teacherID,
			},
			caller: domain.Caller{
				UserID:    callerAdminID,
				Role:      domain.RoleAdmin,
				BranchIDs: []uuid.UUID{sharedBranchID},
			},
			mockBehavior: func(mockReportRepo *mocks.MockReportRepository, mockUserRepo *mocks.MockUserRepository) {
				mockUserRepo.On("GetUserBranchIDs", mock.Anything, teacherID).
					Return([]uuid.UUID{}, nil)
			},
			expectedResult: nil,
			expectedError:  domain.ErrBranchAccessDenied,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockReportRepo := new(mocks.MockReportRepository)
			mockUserRepo := new(mocks.MockUserRepository)
			tt.mockBehavior(mockReportRepo, mockUserRepo)

			uc := NewUseCase(mockReportRepo, mockUserRepo)
			res, err := uc.Execute(context.Background(), tt.caller, tt.req)

			if tt.expectedError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.expectedError)
				require.Nil(t, res)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedResult, res)
			}

			mockReportRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
		})
	}
}
