package domain

import "errors"

var (
	ErrAlreadyExists             = errors.New("already exists")
	ErrNotFound                  = errors.New("not found")
	ErrInvalidInput              = errors.New("invalid input")
	ErrInternal                  = errors.New("internal error")
	ErrPhoneAlreadyExists        = errors.New("phone already exists")
	ErrCannotEditArchived        = errors.New("cannot edit archived entity")
	ErrArchivedReference         = errors.New("referenced entity is archived")
	ErrAlreadyArchived           = errors.New("this entity is already archived")
	ErrAlreadyActive             = errors.New("this entity is already active")
	ErrBranchAccessDenied        = errors.New("branch access denied")
	ErrTeacherScheduleConflict   = errors.New("teacher schedule conflict")
	ErrStudentScheduleConflict   = errors.New("student schedule conflict")
	ErrTeacherHasFutureLessons   = errors.New("teacher has future lessons in branch")
	ErrTeacherHasActiveTemplates = errors.New("teacher has active templates in branch")
)
