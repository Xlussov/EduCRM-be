package list

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/Xlussov/EduCRM-be/internal/controller/http/middleware"
	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/Xlussov/EduCRM-be/pkg/response"
)

type Handler struct {
	usecase *UseCase
}

func NewHandler(uc *UseCase) *Handler {
	return &Handler{usecase: uc}
}

// @Summary List lessons
// @Tags lessons
// @Security BearerAuth
// @Produce json
// @Param from_date query string true "From date" format(date)
// @Param to_date query string true "To date" format(date)
// @Param teacher_id query string false "Teacher ID" format(uuid)
// @Param student_id query string false "Student ID" format(uuid)
// @Param group_id query string false "Group ID" format(uuid)
// @Success 200 {array} LessonResponse "Lessons"
// @Failure 400 {object} response.ErrorResponse "Bad Request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/v1/lessons [get]
func (h *Handler) Handle(c echo.Context) error {
	caller, err := middleware.GetCaller(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error(), nil)
	}

	fromDate := c.QueryParam("from_date")
	toDate := c.QueryParam("to_date")
	if fromDate == "" || toDate == "" {
		return response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "from_date and to_date are required", nil)
	}

	var teacherID *uuid.UUID
	if raw := c.QueryParam("teacher_id"); raw != "" {
		id, err := uuid.Parse(raw)
		if err != nil {
			return response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "invalid teacher_id", nil)
		}
		teacherID = &id
	}

	var studentID *uuid.UUID
	if raw := c.QueryParam("student_id"); raw != "" {
		id, err := uuid.Parse(raw)
		if err != nil {
			return response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "invalid student_id", nil)
		}
		studentID = &id
	}

	var groupID *uuid.UUID
	if raw := c.QueryParam("group_id"); raw != "" {
		id, err := uuid.Parse(raw)
		if err != nil {
			return response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "invalid group_id", nil)
		}
		groupID = &id
	}

	req := Request{
		FromDate:  fromDate,
		ToDate:    toDate,
		TeacherID: teacherID,
		StudentID: studentID,
		GroupID:   groupID,
	}

	res, err := h.usecase.Execute(c.Request().Context(), *caller, req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidInput):
			return response.Error(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
		case errors.Is(err, domain.ErrBranchAccessDenied):
			return response.Error(c, http.StatusForbidden, "BRANCH_ACCESS_DENIED", err.Error(), nil)
		default:
			return response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		}
	}

	return response.Success(c, http.StatusOK, res)
}
