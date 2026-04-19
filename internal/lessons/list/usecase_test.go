package list

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/Xlussov/EduCRM-be/internal/domain/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUseCase_Execute(t *testing.T) {
	branchID := uuid.New()
	teacherID := uuid.New()
	otherTeacherID := uuid.New()
	studentID := uuid.New()
	groupID := uuid.New()
	var nilUUID *uuid.UUID

	fromDateStr := "2026-05-01"
	toDateStr := "2026-05-31"
	fromDate, _ := time.Parse(dateLayout, fromDateStr)
	toDate, _ := time.Parse(dateLayout, toDateStr)

	lessons := []domain.LessonDetails{
		{
			ID:               uuid.New(),
			BranchID:         branchID,
			TeacherID:        teacherID,
			TeacherFirstName: "Alex",
			TeacherLastName:  "Smith",
			SubjectID:        uuid.New(),
			SubjectName:      "Math",
			StudentID:        &studentID,
			StudentFirstName: "John",
			StudentLastName:  "Doe",
			Date:             fromDate,
			StartTime:        mustParseTime(t, "10:00"),
			EndTime:          mustParseTime(t, "11:00"),
			Status:           domain.LessonStatusScheduled,
		},
		{
			ID:               uuid.New(),
			BranchID:         branchID,
			TeacherID:        teacherID,
			TeacherFirstName: "Alex",
			TeacherLastName:  "Smith",
			SubjectID:        uuid.New(),
			SubjectName:      "Science",
			GroupID:          &groupID,
			GroupName:        "Group A",
			Date:             fromDate,
			StartTime:        mustParseTime(t, "12:00"),
			EndTime:          mustParseTime(t, "13:00"),
			Status:           domain.LessonStatusScheduled,
		},
	}

	errDB := errors.New("db error")

	tests := []struct {
		name          string
		caller        domain.Caller
		req           Request
		mockSetup     func(repo *mocks.ScheduleRepository)
		expectedError error
		assertResult  func(t *testing.T, res []LessonResponse)
	}{
		{
			name:   "superadmin_success",
			caller: domain.Caller{Role: domain.RoleSuperadmin},
			req: Request{
				FromDate: fromDateStr,
				ToDate:   toDateStr,
			},
			mockSetup: func(repo *mocks.ScheduleRepository) {
				repo.On(
					"ListLessons",
					mock.Anything,
					fromDate,
					toDate,
					nilUUID,
					nilUUID,
					nilUUID,
					mock.MatchedBy(func(ids []uuid.UUID) bool { return len(ids) == 0 }),
				).Return(lessons, nil).Once()
			},
			assertResult: func(t *testing.T, res []LessonResponse) {
				assert.Len(t, res, 2)
				assert.NotNil(t, res[0].Student)
				assert.Nil(t, res[0].Group)
				assert.Nil(t, res[1].Student)
				assert.NotNil(t, res[1].Group)
			},
		},
		{
			name:   "admin_branch_filter",
			caller: domain.Caller{Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}},
			req: Request{
				FromDate:  fromDateStr,
				ToDate:    toDateStr,
				TeacherID: &teacherID,
			},
			mockSetup: func(repo *mocks.ScheduleRepository) {
				repo.On(
					"ListLessons",
					mock.Anything,
					fromDate,
					toDate,
					&teacherID,
					nilUUID,
					nilUUID,
					[]uuid.UUID{branchID},
				).Return(lessons, nil).Once()
			},
		},
		{
			name:   "teacher_forces_teacher_id",
			caller: domain.Caller{Role: domain.RoleTeacher, UserID: teacherID, BranchIDs: []uuid.UUID{branchID}},
			req: Request{
				FromDate:  fromDateStr,
				ToDate:    toDateStr,
				TeacherID: &otherTeacherID,
			},
			mockSetup: func(repo *mocks.ScheduleRepository) {
				repo.On(
					"ListLessons",
					mock.Anything,
					fromDate,
					toDate,
					mock.MatchedBy(func(id *uuid.UUID) bool { return id != nil && *id == teacherID }),
					nilUUID,
					nilUUID,
					[]uuid.UUID{branchID},
				).Return(lessons, nil).Once()
			},
		},
		{
			name:   "invalid_date",
			caller: domain.Caller{Role: domain.RoleSuperadmin},
			req: Request{
				FromDate: "bad-date",
				ToDate:   toDateStr,
			},
			expectedError: domain.ErrInvalidInput,
		},
		{
			name:   "from_after_to",
			caller: domain.Caller{Role: domain.RoleSuperadmin},
			req: Request{
				FromDate: "2026-06-01",
				ToDate:   "2026-05-01",
			},
			expectedError: domain.ErrInvalidInput,
		},
		{
			name:   "repo_error",
			caller: domain.Caller{Role: domain.RoleSuperadmin},
			req: Request{
				FromDate: fromDateStr,
				ToDate:   toDateStr,
			},
			mockSetup: func(repo *mocks.ScheduleRepository) {
				repo.On(
					"ListLessons",
					mock.Anything,
					fromDate,
					toDate,
					nilUUID,
					nilUUID,
					nilUUID,
					mock.MatchedBy(func(ids []uuid.UUID) bool { return len(ids) == 0 }),
				).Return(nil, errDB).Once()
			},
			expectedError: errDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.ScheduleRepository)
			if tt.mockSetup != nil {
				tt.mockSetup(repo)
			}

			uc := NewUseCase(repo)
			res, err := uc.Execute(context.Background(), tt.caller, tt.req)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
				return
			}
			assert.NoError(t, err)

			if tt.assertResult != nil {
				tt.assertResult(t, res)
			}

			repo.AssertExpectations(t)
		})
	}
}

func mustParseTime(t *testing.T, value string) time.Time {
	parsed, err := time.Parse(timeLayout, value)
	if err != nil {
		t.Fatalf("failed to parse time: %v", err)
	}
	return parsed
}
