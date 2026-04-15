package get

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
	subjectID := uuid.New()
	branchID := uuid.New()
	now := time.Now()

	expectedSubject := &domain.Subject{
		ID:          subjectID,
		BranchID:    branchID,
		Name:        "Test Subject",
		Description: "A description",
		Status:      domain.StatusActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	tests := []struct {
		name      string
		req       Request
		mockSetup func(repo *mocks.SubjectRepository)
		wantRes   *Response
		wantErr   bool
	}{
		{
			name: "Success: Repository successfully returns the subject",
			req:  Request{ID: subjectID},
			mockSetup: func(repo *mocks.SubjectRepository) {
				repo.On("GetByID", mock.Anything, subjectID).Return(expectedSubject, nil)
			},
			wantRes: &Response{
				Subject: SubjectResponse{
					ID:          subjectID.String(),
					BranchID:    branchID.String(),
					Name:        "Test Subject",
					Description: "A description",
					Status:      domain.StatusActive,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			wantErr: false,
		},
		{
			name: "Repo Error: Repository returns an error",
			req:  Request{ID: subjectID},
			mockSetup: func(repo *mocks.SubjectRepository) {
				repo.On("GetByID", mock.Anything, subjectID).Return(nil, errors.New("db error"))
			},
			wantRes: nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.SubjectRepository)
			if tt.mockSetup != nil {
				tt.mockSetup(repo)
			}

			uc := NewUseCase(repo)
			res, err := uc.Execute(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantRes, res)
			}

			repo.AssertExpectations(t)
		})
	}
}
