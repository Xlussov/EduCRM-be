package mark

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

// @Summary      Mark lesson attendance
// @Description  Marks attendance for the lesson students
// @Tags         attendance
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "Lesson ID"  format(uuid)
// @Param        request body Request true "Attendance payload"
// @Success      200  {object}  Response
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
// @Failure      403  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /api/v1/attendance/{id} [put]
func (h *Handler) Handle(c echo.Context) error {
	caller, err := middleware.GetCaller(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error(), nil)
	}

	lessonID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "invalid lesson ID", nil)
	}

	var req Request
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "invalid request body", nil)
	}
	if req.Attendance == nil {
		return response.Error(c, http.StatusBadRequest, "VALIDATION_FAILED", "attendance is required", nil)
	}

	res, err := h.usecase.Execute(c.Request().Context(), *caller, lessonID, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrAttendanceRequired), errors.Is(err, ErrIsPresentRequired):
			return response.Error(c, http.StatusBadRequest, "VALIDATION_FAILED", err.Error(), nil)
		case errors.Is(err, ErrLessonCancelled):
			return response.Error(c, http.StatusBadRequest, "LESSON_IS_CANCELLED", err.Error(), nil)
		case errors.Is(err, ErrLessonInFuture):
			return response.Error(c, http.StatusBadRequest, "LESSON_IN_FUTURE", err.Error(), nil)
		case errors.Is(err, ErrStudentNotInLesson):
			return response.Error(c, http.StatusBadRequest, "STUDENT_NOT_IN_LESSON", err.Error(), nil)
		case errors.Is(err, ErrDuplicateStudentID):
			return response.Error(c, http.StatusBadRequest, "DUPLICATE_STUDENT", err.Error(), nil)
		case errors.Is(err, domain.ErrBranchAccessDenied):
			return response.Error(c, http.StatusForbidden, "BRANCH_ACCESS_DENIED", err.Error(), nil)
		case errors.Is(err, domain.ErrNotFound):
			return response.Error(c, http.StatusNotFound, "NOT_FOUND", "lesson not found", nil)
		default:
			return response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		}
	}

	return response.Success(c, http.StatusOK, res)
}
