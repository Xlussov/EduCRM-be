package list

import (
	"context"
	"testing"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/Xlussov/EduCRM-be/internal/domain/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUseCase_Execute(t *testing.T) {
	branch1 := uuid.New()
	branch2 := uuid.New()
	userID := uuid.New()

	callerSuper := domain.Caller{UserID: userID, Role: domain.RoleSuperadmin}
	callerAdmin := domain.Caller{UserID: userID, Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branch1}}
	callerAdminNoAccess := domain.Caller{UserID: userID, Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branch2}}
	callerTeacher := domain.Caller{UserID: userID, Role: domain.RoleTeacher, BranchIDs: []uuid.UUID{branch1}}
	callerTeacherNoAccess := domain.Caller{UserID: userID, Role: domain.RoleTeacher, BranchIDs: []uuid.UUID{branch2}}

	tests := []struct {
		name        string
		caller      domain.Caller
		req         Request
		setupMocks  func(mockGR *mocks.GroupRepository)
		expectedErr error
	}{
		{
			name:   "Success_SUPERADMIN",
			caller: callerSuper,
			req:    Request{BranchID: branch1},
			setupMocks: func(mockGR *mocks.GroupRepository) {
				mockGR.On("GetByBranchID", mock.Anything, branch1, (*domain.EntityStatus)(nil)).Return([]*domain.GroupWithCount{
					{
						Group: domain.Group{
							ID:     uuid.New(),
							Name:   "Test Group",
							Status: domain.StatusActive,
						},
						StudentsCount: 5,
					},
				}, nil).Once()
			},
			expectedErr: nil,
		},
		{
			name:   "Success_ADMIN_HasAccess",
			caller: callerAdmin,
			req:    Request{BranchID: branch1},
			setupMocks: func(mockGR *mocks.GroupRepository) {
				mockGR.On("GetByBranchID", mock.Anything, branch1, (*domain.EntityStatus)(nil)).Return([]*domain.GroupWithCount{}, nil).Once()
			},
			expectedErr: nil,
		},
		{
			name:   "Success_TEACHER_HasAccess",
			caller: callerTeacher,
			req:    Request{BranchID: branch1},
			setupMocks: func(mockGR *mocks.GroupRepository) {
				mockGR.On("GetByBranchIDAndTeacherID", mock.Anything, branch1, userID, (*domain.EntityStatus)(nil)).Return([]*domain.GroupWithCount{}, nil).Once()
			},
			expectedErr: nil,
		},
		{
			name:        "Forbidden_ADMIN_NoAccess",
			caller:      callerAdminNoAccess,
			req:         Request{BranchID: branch1},
			expectedErr: domain.ErrBranchAccessDenied,
		},
		{
			name:        "Forbidden_TEACHER_NoAccess",
			caller:      callerTeacherNoAccess,
			req:         Request{BranchID: branch1},
			expectedErr: domain.ErrBranchAccessDenied,
		},
		{
			name:        "Error_MissingBranchID",
			caller:      callerSuper,
			req:         Request{BranchID: uuid.Nil},
			setupMocks:  func(mockGR *mocks.GroupRepository) {},
			expectedErr: ErrBranchIDRequired,
		},
		{
			name:        "Error_InvalidStatus",
			caller:      callerSuper,
			req:         Request{BranchID: branch1, Status: "BROKEN"},
			setupMocks:  func(mockGR *mocks.GroupRepository) {},
			expectedErr: domain.ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGR := new(mocks.GroupRepository)
			if tt.setupMocks != nil {
				tt.setupMocks(mockGR)
			}

			uc := NewUseCase(mockGR)
			_, err := uc.Execute(context.Background(), tt.caller, tt.req)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
			}

			mockGR.AssertExpectations(t)
		})
	}
}
