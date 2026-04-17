package create_template

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Xlussov/EduCRM-be/internal/domain"
)

type UseCase struct {
	scheduleRepo domain.ScheduleRepository
	groupRepo    domain.GroupRepository
}

func NewUseCase(sr domain.ScheduleRepository, gr domain.GroupRepository) *UseCase {
	return &UseCase{
		scheduleRepo: sr,
		groupRepo:    gr,
	}
}

func (uc *UseCase) Execute(ctx context.Context, caller domain.Caller, req Request) (Response, error) {
	if (req.StudentID == nil && req.GroupID == nil) || (req.StudentID != nil && req.GroupID != nil) {
		return Response{}, domain.ErrInvalidInput
	}

	if domain.RequiresBranchAccess(caller.Role) && !domain.HasBranchAccess(caller.BranchIDs, req.BranchID) {
		return Response{}, domain.ErrBranchAccessDenied
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return Response{}, domain.ErrInvalidInput
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return Response{}, domain.ErrInvalidInput
	}

	if !startDate.Before(endDate) && !startDate.Equal(endDate) {
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

	var studentIDs []uuid.UUID
	if req.GroupID != nil {
		ids, err := uc.groupRepo.GetActiveStudentIDs(ctx, *req.GroupID)
		if err != nil {
			return Response{}, err
		}
		studentIDs = ids
	} else if req.StudentID != nil {
		studentIDs = []uuid.UUID{*req.StudentID}
	}

	template := &domain.Template{
		BranchID:  req.BranchID,
		TeacherID: req.TeacherID,
		SubjectID: req.SubjectID,
		StudentID: req.StudentID,
		GroupID:   req.GroupID,
		DayOfWeek: req.DayOfWeek,
		StartTime: start,
		EndTime:   end,
		StartDate: startDate,
		EndDate:   endDate,
		IsActive:  true,
	}

	if err := uc.scheduleRepo.CreateTemplate(ctx, template); err != nil {
		return Response{}, err
	}

	var conflicts []string
	var lessonsToCreate []domain.Lesson

	currentDate := startDate
	for !currentDate.After(endDate) {
		if int(currentDate.Weekday()) == req.DayOfWeek {
			hasConflict := false

			tConflict, err := uc.scheduleRepo.CheckTeacherConflict(ctx, req.TeacherID, currentDate, start, end)
			if err != nil {
				return Response{}, err
			}

			if tConflict {
				hasConflict = true
			} else {
				for _, sid := range studentIDs {
					sConflict, err := uc.scheduleRepo.CheckStudentConflict(ctx, sid, currentDate, start, end)
					if err != nil {
						return Response{}, err
					}
					if sConflict {
						hasConflict = true
						break
					}
				}
			}

			if hasConflict {
				conflicts = append(conflicts, currentDate.Format("2006-01-02"))
			} else {
				lessonsToCreate = append(lessonsToCreate, domain.Lesson{
					BranchID:   req.BranchID,
					TemplateID: &template.ID,
					TeacherID:  req.TeacherID,
					SubjectID:  req.SubjectID,
					StudentID:  req.StudentID,
					GroupID:    req.GroupID,
					Date:       currentDate,
					StartTime:  start,
					EndTime:    end,
					Status:     domain.LessonStatusScheduled,
				})
			}
		}

		currentDate = currentDate.AddDate(0, 0, 1)
	}

	if len(lessonsToCreate) > 0 {
		if err := uc.scheduleRepo.BulkCreateLessons(ctx, lessonsToCreate); err != nil {
			return Response{}, err
		}
	}

	if conflicts == nil {
		conflicts = make([]string, 0)
	}

	return Response{
		TemplateID:          template.ID,
		CreatedLessonsCount: len(lessonsToCreate),
		Conflicts:           conflicts,
	}, nil
}
