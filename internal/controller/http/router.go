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
	AuthLogin   echo.HandlerFunc
	AuthRefresh echo.HandlerFunc
	AuthLogout  echo.HandlerFunc
	AuthMe      echo.HandlerFunc

	UsersAdminsArchive   echo.HandlerFunc
	UsersAdminsGet       echo.HandlerFunc
	UsersAdminsList      echo.HandlerFunc
	UsersAdminsCreate    echo.HandlerFunc
	UsersAdminsUnarchive echo.HandlerFunc
	UsersAdminsUpdate    echo.HandlerFunc
	UsersTeachersCreate  echo.HandlerFunc

	BranchesCreate    echo.HandlerFunc
	BranchesArchive   echo.HandlerFunc
	BranchesUnarchive echo.HandlerFunc
	BranchesList      echo.HandlerFunc
	BranchesGet       echo.HandlerFunc
	BranchesUpdate    echo.HandlerFunc

	SubjectsCreate    echo.HandlerFunc
	SubjectsArchive   echo.HandlerFunc
	SubjectsUnarchive echo.HandlerFunc
	SubjectsList      echo.HandlerFunc
	SubjectsUpdate    echo.HandlerFunc

	StudentsCreate    echo.HandlerFunc
	StudentsArchive   echo.HandlerFunc
	StudentsUnarchive echo.HandlerFunc
	StudentsList      echo.HandlerFunc
	StudentsGet       echo.HandlerFunc
	StudentsUpdate    echo.HandlerFunc

	GroupsCreate        echo.HandlerFunc
	GroupsList          echo.HandlerFunc
	GroupsGet           echo.HandlerFunc
	GroupsUpdate        echo.HandlerFunc
	GroupsAddStudents   echo.HandlerFunc
	GroupsRemoveStudent echo.HandlerFunc
	GroupsArchive       echo.HandlerFunc
	GroupsUnarchive     echo.HandlerFunc

	PlansCreate  echo.HandlerFunc
	PlansList    echo.HandlerFunc
	PlansArchive echo.HandlerFunc

	SubscriptionsCreate echo.HandlerFunc
	SubscriptionsList   echo.HandlerFunc
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

		protectedAuthGroup := protected.Group("/auth")
		protectedAuthGroup.GET("/me", h.AuthMe)
		protectedAuthGroup.POST("/logout", h.AuthLogout)

		usersGroup := protected.Group("/users")

		adminsGroup := usersGroup.Group("/admins")
		adminsGroup.POST("", h.UsersAdminsCreate)
		adminsGroup.GET("", h.UsersAdminsList)
		adminsGroup.GET("/:id", h.UsersAdminsGet)
		adminsGroup.PUT("/:id", h.UsersAdminsUpdate)
		adminsGroup.PATCH("/:id/archive", h.UsersAdminsArchive)
		adminsGroup.PATCH("/:id/unarchive", h.UsersAdminsUnarchive)

		teachersGroup := usersGroup.Group("/teachers")
		teachersGroup.POST("", h.UsersTeachersCreate)

		branchesGroup := protected.Group("/branches")
		branchesGroup.POST("", h.BranchesCreate)
		branchesGroup.GET("", h.BranchesList)
		branchesGroup.GET("/:id", h.BranchesGet)
		branchesGroup.PUT("/:id", h.BranchesUpdate)
		branchesGroup.PATCH("/:id/archive", h.BranchesArchive)
		branchesGroup.PATCH("/:id/unarchive", h.BranchesUnarchive)

		subjectsGroup := protected.Group("/subjects")
		subjectsGroup.POST("", h.SubjectsCreate)
		subjectsGroup.GET("", h.SubjectsList)
		subjectsGroup.PUT("/:id", h.SubjectsUpdate)
		subjectsGroup.PATCH("/:id/archive", h.SubjectsArchive)
		subjectsGroup.PATCH("/:id/unarchive", h.SubjectsUnarchive)

		studentsGroup := protected.Group("/students")
		studentsGroup.POST("", h.StudentsCreate)
		studentsGroup.GET("", h.StudentsList)
		studentsGroup.GET("/:id", h.StudentsGet)
		studentsGroup.PUT("/:id", h.StudentsUpdate)
		studentsGroup.PATCH("/:id/archive", h.StudentsArchive)
		studentsGroup.PATCH("/:id/unarchive", h.StudentsUnarchive)
		studentsGroup.POST("/:id/subscriptions", h.SubscriptionsCreate)
		studentsGroup.GET("/:id/subscriptions", h.SubscriptionsList)

		groupsGroup := protected.Group("/groups")
		groupsGroup.POST("", h.GroupsCreate)
		groupsGroup.GET("", h.GroupsList)
		groupsGroup.GET("/:id", h.GroupsGet)
		groupsGroup.PUT("/:id", h.GroupsUpdate)
		groupsGroup.POST("/:id/students", h.GroupsAddStudents)
		groupsGroup.DELETE("/:id/students/:student_id", h.GroupsRemoveStudent)
		groupsGroup.PATCH("/:id/archive", h.GroupsArchive)
		groupsGroup.PATCH("/:id/unarchive", h.GroupsUnarchive)

		plansGroup := protected.Group("/plans")
		plansGroup.POST("", h.PlansCreate)
		plansGroup.GET("", h.PlansList)
		plansGroup.PATCH("/:id/archive", h.PlansArchive)
	}
}
