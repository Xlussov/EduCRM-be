package get

import (
	"context"
	"fmt"

	"github.com/Xlussov/EduCRM-be/internal/domain"
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

func (uc *UseCase) Execute(ctx context.Context, caller domain.Caller, req Request) (Response, error) {
	lesson, err := uc.scheduleRepo.GetLessonByID(ctx, req.ID)
	if err != nil {
		return Response{}, fmt.Errorf("get lesson: %w", err)
	}

	if domain.RequiresBranchAccess(caller.Role) && !domain.HasBranchAccess(caller.BranchIDs, lesson.BranchID) {
		return Response{}, domain.ErrBranchAccessDenied
	}

	if caller.Role == domain.RoleTeacher && lesson.TeacherID != caller.UserID {
		return Response{}, domain.ErrBranchAccessDenied
	}

	return Response{
		ID:         lesson.ID,
		BranchID:   lesson.BranchID,
		TemplateID: lesson.TemplateID,
		TeacherID:  lesson.TeacherID,
		SubjectID:  lesson.SubjectID,
		StudentID:  lesson.StudentID,
		GroupID:    lesson.GroupID,
		Date:       lesson.Date.Format(dateLayout),
		StartTime:  lesson.StartTime.Format(timeLayout),
		EndTime:    lesson.EndTime.Format(timeLayout),
		Status:     string(lesson.Status),
	}, nil
}
