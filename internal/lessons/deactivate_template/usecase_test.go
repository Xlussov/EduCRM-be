package deactivate_template

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
	templateID := uuid.New()
	branchID := uuid.New()
	otherBranchID := uuid.New()
	errDB := errors.New("db error")

	activeTemplate := &domain.Template{ID: templateID, BranchID: branchID, IsActive: true}
	inactiveTemplate := &domain.Template{ID: templateID, BranchID: branchID, IsActive: false}

	tests := []struct {
		name        string
		caller      domain.Caller
		mockSetup   func(repo *mocks.ScheduleRepository)
		expectedErr error
		assertions  func(t *testing.T, repo *mocks.ScheduleRepository)
	}{
		{
			name:   "superadmin_success",
			caller: domain.Caller{Role: domain.RoleSuperadmin},
			mockSetup: func(repo *mocks.ScheduleRepository) {
				repo.On("GetTemplateByID", mock.Anything, templateID).Return(activeTemplate, nil).Once()
				repo.On("DeactivateTemplate", mock.Anything, templateID).Return(nil).Once()
				repo.On("CancelFutureLessonsByTemplate", mock.Anything, templateID).Return(nil).Once()
			},
		},
		{
			name:   "admin_success",
			caller: domain.Caller{Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}},
			mockSetup: func(repo *mocks.ScheduleRepository) {
				repo.On("GetTemplateByID", mock.Anything, templateID).Return(activeTemplate, nil).Once()
				repo.On("DeactivateTemplate", mock.Anything, templateID).Return(nil).Once()
				repo.On("CancelFutureLessonsByTemplate", mock.Anything, templateID).Return(nil).Once()
			},
		},
		{
			name:        "admin_access_denied",
			caller:      domain.Caller{Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{otherBranchID}},
			expectedErr: domain.ErrBranchAccessDenied,
			mockSetup: func(repo *mocks.ScheduleRepository) {
				repo.On("GetTemplateByID", mock.Anything, templateID).Return(activeTemplate, nil).Once()
			},
			assertions: func(t *testing.T, repo *mocks.ScheduleRepository) {
				repo.AssertNotCalled(t, "DeactivateTemplate", mock.Anything, templateID)
				repo.AssertNotCalled(t, "CancelFutureLessonsByTemplate", mock.Anything, templateID)
			},
		},
		{
			name:        "template_not_found",
			caller:      domain.Caller{Role: domain.RoleSuperadmin},
			expectedErr: domain.ErrNotFound,
			mockSetup: func(repo *mocks.ScheduleRepository) {
				repo.On("GetTemplateByID", mock.Anything, templateID).Return(nil, domain.ErrNotFound).Once()
			},
		},
		{
			name:        "template_not_active",
			caller:      domain.Caller{Role: domain.RoleSuperadmin},
			expectedErr: domain.ErrTemplateNotActive,
			mockSetup: func(repo *mocks.ScheduleRepository) {
				repo.On("GetTemplateByID", mock.Anything, templateID).Return(inactiveTemplate, nil).Once()
			},
		},
		{
			name:        "deactivate_error",
			caller:      domain.Caller{Role: domain.RoleSuperadmin},
			expectedErr: errDB,
			mockSetup: func(repo *mocks.ScheduleRepository) {
				repo.On("GetTemplateByID", mock.Anything, templateID).Return(activeTemplate, nil).Once()
				repo.On("DeactivateTemplate", mock.Anything, templateID).Return(errDB).Once()
			},
		},
		{
			name:        "cancel_future_lessons_error",
			caller:      domain.Caller{Role: domain.RoleSuperadmin},
			expectedErr: errDB,
			mockSetup: func(repo *mocks.ScheduleRepository) {
				repo.On("GetTemplateByID", mock.Anything, templateID).Return(activeTemplate, nil).Once()
				repo.On("DeactivateTemplate", mock.Anything, templateID).Return(nil).Once()
				repo.On("CancelFutureLessonsByTemplate", mock.Anything, templateID).Return(errDB).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.ScheduleRepository)
			if tt.mockSetup != nil {
				tt.mockSetup(repo)
			}

			uc := NewUseCase(repo, &mocks.MockTxManager{})
			res, err := uc.Execute(context.Background(), tt.caller, templateID)

			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "success", res.Message)
			}

			if tt.assertions != nil {
				tt.assertions(t, repo)
			}

			repo.AssertExpectations(t)
		})
	}
}
