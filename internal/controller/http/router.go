package http

import (
	"net/http"

	mw "github.com/Xlussov/EduCRM-be/internal/controller/http/middleware"
	"github.com/Xlussov/EduCRM-be/pkg/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Handlers struct {
	AuthLogin       echo.HandlerFunc
	AuthRefresh     echo.HandlerFunc
	UsersAdmins     echo.HandlerFunc
	UsersTeachers   echo.HandlerFunc
	BranchesCreate  echo.HandlerFunc
	BranchesArchive echo.HandlerFunc
	BranchesList    echo.HandlerFunc
	BranchesGet     echo.HandlerFunc
	BranchesUpdate  echo.HandlerFunc
	SubjectsCreate  echo.HandlerFunc
	SubjectsArchive echo.HandlerFunc
	SubjectsList    echo.HandlerFunc
	SubjectsUpdate  echo.HandlerFunc
	StudentsCreate  echo.HandlerFunc
	StudentsArchive echo.HandlerFunc
	StudentsList    echo.HandlerFunc
	StudentsGet     echo.HandlerFunc
	StudentsUpdate  echo.HandlerFunc
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Error(msg string)
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Errorf(format string, args ...any)
}

func Init(log Logger, cfg *config.Config, e *echo.Echo, h Handlers) {
	e.Use(mw.RequestLogger(log))
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.CORS())

	e.GET("/docs", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	v1 := e.Group("/api/v1")
	{
		// Public routes (Auth)
		authGroup := v1.Group("/auth")
		authGroup.POST("/login", h.AuthLogin)
		authGroup.POST("/refresh", h.AuthRefresh)

		// Protected routes
		protected := v1.Group("")
		protected.Use(mw.JWT(cfg.JWTSecret))

		usersGroup := protected.Group("/users")
		usersGroup.POST("/admins", h.UsersAdmins)
		usersGroup.POST("/teachers", h.UsersTeachers)

		branchesGroup := protected.Group("/branches")
		branchesGroup.POST("", h.BranchesCreate)
		branchesGroup.GET("", h.BranchesList)
		branchesGroup.GET("/:id", h.BranchesGet)
		branchesGroup.PUT("/:id", h.BranchesUpdate)
		branchesGroup.PATCH("/:id/archive", h.BranchesArchive)
		subjectsGroup := protected.Group("/subjects")
		subjectsGroup.POST("", h.SubjectsCreate)
		subjectsGroup.GET("", h.SubjectsList)
		subjectsGroup.PUT("/:id", h.SubjectsUpdate)
		subjectsGroup.PATCH("/:id/archive", h.SubjectsArchive)

		studentsGroup := protected.Group("/students")
		studentsGroup.POST("", h.StudentsCreate)
		studentsGroup.GET("", h.StudentsList)
		studentsGroup.GET("/:id", h.StudentsGet)
		studentsGroup.PUT("/:id", h.StudentsUpdate)
		studentsGroup.PATCH("/:id/archive", h.StudentsArchive)
	}
}
