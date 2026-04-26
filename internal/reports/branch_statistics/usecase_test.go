package branch_statistics

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
	branchID := uuid.New()
	callerID := uuid.New()
	now := time.Now()

	tests := []struct {
		name           string
		req            Request
		caller         domain.Caller
		mockBehavior   func(mockReportRepo *mocks.MockReportRepository)
		expectedResult *Response
		expectedError  error
	}{
		{
			name: "Success - Admin with branch access",
			req: Request{
				BranchID:  branchID,
				StartDate: &now,
				EndDate:   &now,
			},
			caller: domain.Caller{
				UserID:    callerID,
				Role:      domain.RoleAdmin,
				BranchIDs: []uuid.UUID{branchID},
			},
			mockBehavior: func(mockReportRepo *mocks.MockReportRepository) {
				mockReportRepo.On("GetBranchStatistics", mock.Anything, branchID, &now, &now).
					Return(&domain.BranchStatisticsReport{
						ActiveStudents:       10,
						CompletedLessons:     5,
						CancelledLessons:     1,
						AttendancePercentage: 95.5,
					}, nil)
			},
			expectedResult: &Response{
				ActiveStudents:       10,
				CompletedLessons:     5,
				CancelledLessons:     1,
				AttendancePercentage: 95.5,
			},
			expectedError: nil,
		},
		{
			name: "Error - Admin without branch access",
			req: Request{
				BranchID:  branchID,
				StartDate: &now,
				EndDate:   &now,
			},
			caller: domain.Caller{
				UserID:    callerID,
				Role:      domain.RoleAdmin,
				BranchIDs: []uuid.UUID{uuid.New()}, // Different branch
			},
			mockBehavior: func(mockReportRepo *mocks.MockReportRepository) {
				// Should not call repo
			},
			expectedResult: nil,
			expectedError:  domain.ErrBranchAccessDenied,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockReportRepo := new(mocks.MockReportRepository)
			tt.mockBehavior(mockReportRepo)

			uc := NewUseCase(mockReportRepo)
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
		})
	}
}
