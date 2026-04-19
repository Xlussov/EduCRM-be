package update

import (
	"context"
	"time"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
)

type UseCase struct {
	scheduleRepo domain.ScheduleRepository
	groupRepo    domain.GroupRepository
	userRepo     domain.UserRepository
}

func NewUseCase(sr domain.ScheduleRepository, gr domain.GroupRepository, ur domain.UserRepository) *UseCase {
	return &UseCase{scheduleRepo: sr, groupRepo: gr, userRepo: ur}
}

func (uc *UseCase) Execute(ctx context.Context, caller domain.Caller, lessonID uuid.UUID, req Request) (Response, error) {
	lesson, err := uc.scheduleRepo.GetLessonByID(ctx, lessonID)
	if err != nil {
		return Response{}, err
	}

	if domain.RequiresBranchAccess(caller.Role) && !domain.HasBranchAccess(caller.BranchIDs, lesson.BranchID) {
		return Response{}, domain.ErrBranchAccessDenied
	}

	if lesson.Status != domain.LessonStatusScheduled {
		return Response{}, domain.ErrLessonNotScheduled
	}

	teacherInBranch, err := uc.userRepo.CheckTeacherInBranch(ctx, req.TeacherID, lesson.BranchID)
	if err != nil {
		return Response{}, err
	}
	if !teacherInBranch {
		return Response{}, domain.ErrTeacherNotInBranch
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return Response{}, domain.ErrInvalidInput
	}

	start, err := time.Parse("15:04", req.StartTime)
	if err != nil {
		return Response{}, domain.ErrInvalidInput
	}

	end, err := time.Parse("15:04", req.EndTime)
	if err != nil {
		return Response{}, domain.ErrInvalidInput
	}

	if !start.Before(end) {
		return Response{}, domain.ErrInvalidInput
	}

	if lesson.StudentID != nil {
		hasStudentConflict, err := uc.scheduleRepo.CheckStudentConflictExcludingLesson(ctx, *lesson.StudentID, date, start, end, lessonID)
		if err != nil {
			return Response{}, err
		}
		if hasStudentConflict {
			return Response{}, domain.ErrStudentScheduleConflict
		}
	}

	if lesson.GroupID != nil {
		studentIDs, err := uc.groupRepo.GetActiveStudentIDs(ctx, *lesson.GroupID)
		if err != nil {
			return Response{}, err
		}

		for _, studentID := range studentIDs {
			hasConflict, err := uc.scheduleRepo.CheckStudentConflictExcludingLesson(ctx, studentID, date, start, end, lessonID)
			if err != nil {
				return Response{}, err
			}
			if hasConflict {
				return Response{}, domain.ErrStudentScheduleConflict
			}
		}
	}

	hasTeacherConflict, err := uc.scheduleRepo.CheckTeacherConflictExcludingLesson(ctx, req.TeacherID, date, start, end, lessonID)
	if err != nil {
		return Response{}, err
	}
	if hasTeacherConflict {
		return Response{}, domain.ErrTeacherScheduleConflict
	}

	lesson.Date = date
	lesson.StartTime = start
	lesson.EndTime = end
	lesson.TeacherID = req.TeacherID
	lesson.SubjectID = req.SubjectID

	if err := uc.scheduleRepo.UpdateLesson(ctx, lesson); err != nil {
		return Response{}, err
	}

	return Response{
		ID:     lesson.ID,
		Status: string(lesson.Status),
	}, nil
}
