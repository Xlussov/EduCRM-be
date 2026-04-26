package student_attendance

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
	studentID := uuid.New()
	callerID := uuid.New()
	branchID := uuid.New()
	now := time.Now()

	items := []domain.StudentAttendanceReportItem{
		{
			Date:        now,
			Time:        now,
			SubjectName: "Math",
			IsPresent:   true,
			Notes:       "",
		},
		{
			Date:        now,
			Time:        now,
			SubjectName: "Math",
			IsPresent:   false,
			Notes:       "Sick",
		},
	}

	tests := []struct {
		name           string
		req            Request
		caller         domain.Caller
		mockBehavior   func(mockReportRepo *mocks.MockReportRepository, mockStudentRepo *mocks.MockStudentRepository)
		expectedResult *Response
		expectedError  error
	}{
		{
			name: "Success - Admin with branch access",
			req: Request{
				StudentID: studentID,
			},
			caller: domain.Caller{
				UserID:    callerID,
				Role:      domain.RoleAdmin,
				BranchIDs: []uuid.UUID{branchID},
			},
			mockBehavior: func(mockReportRepo *mocks.MockReportRepository, mockStudentRepo *mocks.MockStudentRepository) {
				mockStudentRepo.On("GetByID", mock.Anything, studentID).
					Return(&domain.Student{ID: studentID, BranchID: branchID}, nil)

				mockReportRepo.On("GetStudentAttendanceHistory", mock.Anything, studentID, (*time.Time)(nil), (*time.Time)(nil), (*uuid.UUID)(nil)).
					Return(items, nil)
			},
			expectedResult: &Response{
				Summary: Summary{
					TotalLessons:         2,
					Attended:             1,
					Missed:               1,
					AttendancePercentage: 50.0,
				},
				Items: []ReportItem{
					{
						Date:        now,
						Time:        now,
						SubjectName: "Math",
						IsPresent:   true,
						Notes:       "",
					},
					{
						Date:        now,
						Time:        now,
						SubjectName: "Math",
						IsPresent:   false,
						Notes:       "Sick",
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "Error - Admin without branch access",
			req: Request{
				StudentID: studentID,
			},
			caller: domain.Caller{
				UserID:    callerID,
				Role:      domain.RoleAdmin,
				BranchIDs: []uuid.UUID{uuid.New()}, // Different branch
			},
			mockBehavior: func(mockReportRepo *mocks.MockReportRepository, mockStudentRepo *mocks.MockStudentRepository) {
				mockStudentRepo.On("GetByID", mock.Anything, studentID).
					Return(&domain.Student{ID: studentID, BranchID: branchID}, nil)
			},
			expectedResult: nil,
			expectedError:  domain.ErrBranchAccessDenied,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockReportRepo := new(mocks.MockReportRepository)
			mockStudentRepo := new(mocks.MockStudentRepository)
			tt.mockBehavior(mockReportRepo, mockStudentRepo)

			uc := NewUseCase(mockReportRepo, mockStudentRepo)
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
			mockStudentRepo.AssertExpectations(t)
		})
	}
}
