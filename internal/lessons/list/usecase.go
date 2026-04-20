package list

import (
	"context"
	"time"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
)

const (
	dateLayout = "2006-01-02"
	timeLayout = "15:04"
)

type UseCase struct {
	scheduleRepo domain.ScheduleRepository
}

func NewUseCase(sr domain.ScheduleRepository) *UseCase {
	return &UseCase{scheduleRepo: sr}
}

func (uc *UseCase) Execute(ctx context.Context, caller domain.Caller, req Request) ([]LessonResponse, error) {
	if req.FromDate == "" || req.ToDate == "" {
		return nil, domain.ErrInvalidInput
	}

	fromDate, err := time.Parse(dateLayout, req.FromDate)
	if err != nil {
		return nil, domain.ErrInvalidInput
	}

	toDate, err := time.Parse(dateLayout, req.ToDate)
	if err != nil {
		return nil, domain.ErrInvalidInput
	}

	if fromDate.After(toDate) {
		return nil, domain.ErrInvalidInput
	}

	teacherID := req.TeacherID
	if caller.Role == domain.RoleTeacher {
		teacherID = &caller.UserID
	}

	var branchFilter []uuid.UUID
	if domain.RequiresBranchAccess(caller.Role) {
		branchFilter = caller.BranchIDs
	}

	lessons, err := uc.scheduleRepo.ListLessons(ctx, fromDate, toDate, teacherID, req.StudentID, req.GroupID, branchFilter)
	if err != nil {
		return nil, err
	}

	res := make([]LessonResponse, 0, len(lessons))
	for _, l := range lessons {
		var student *StudentRef
		if l.StudentID != nil {
			student = &StudentRef{ID: *l.StudentID, FirstName: l.StudentFirstName, LastName: l.StudentLastName}
		}

		var group *GroupRef
		if l.GroupID != nil {
			group = &GroupRef{ID: *l.GroupID, Name: l.GroupName}
		}

		res = append(res, LessonResponse{
			ID:         l.ID,
			TemplateID: l.TemplateID,
			Date:       l.Date.Format(dateLayout),
			StartTime:  l.StartTime.Format(timeLayout),
			EndTime:    l.EndTime.Format(timeLayout),
			Status:     string(l.Status),
			Teacher:    TeacherRef{ID: l.TeacherID, FirstName: l.TeacherFirstName, LastName: l.TeacherLastName},
			Subject:    SubjectRef{ID: l.SubjectID, Name: l.SubjectName},
			Student:    student,
			Group:      group,
		})
	}

	return res, nil
}
