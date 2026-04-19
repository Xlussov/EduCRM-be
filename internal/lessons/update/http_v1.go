package update

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

// @Summary      Update lesson
// @Description  Updates lesson date, time, teacher, and subject
// @Tags         lessons
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "Lesson ID"  format(uuid)
// @Param        request body Request true "Lesson update payload"
// @Success      200  {object}  Response
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
// @Failure      403  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      409  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /api/v1/lessons/{id} [patch]
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

	if err := c.Validate(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "VALIDATION_FAILED", err.Error(), nil)
	}

	res, err := h.usecase.Execute(c.Request().Context(), *caller, lessonID, req)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidInput) {
			return response.Error(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
		}
		if errors.Is(err, domain.ErrTeacherNotInBranch) {
			return response.Error(c, http.StatusBadRequest, "TEACHER_NOT_IN_BRANCH", err.Error(), nil)
		}
		if errors.Is(err, domain.ErrBranchAccessDenied) {
			return response.Error(c, http.StatusForbidden, "FORBIDDEN", err.Error(), nil)
		}
		if errors.Is(err, domain.ErrNotFound) {
			return response.Error(c, http.StatusNotFound, "NOT_FOUND", "lesson not found", nil)
		}
		if errors.Is(err, domain.ErrLessonNotScheduled) {
			return response.Error(c, http.StatusConflict, "LESSON_NOT_SCHEDULED", err.Error(), nil)
		}
		if errors.Is(err, domain.ErrTeacherScheduleConflict) || errors.Is(err, domain.ErrStudentScheduleConflict) {
			return response.Error(c, http.StatusConflict, "CONFLICT", err.Error(), nil)
		}
		return response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
	}

	return response.Success(c, http.StatusOK, res)
}
