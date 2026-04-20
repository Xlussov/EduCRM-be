package http

import (
	"net/http"

	mw "github.com/Xlussov/EduCRM-be/internal/controller/http/middleware"
	"github.com/Xlussov/EduCRM-be/internal/domain"
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

	UsersAdminsArchive     echo.HandlerFunc
	UsersAdminsGet         echo.HandlerFunc
	UsersAdminsList        echo.HandlerFunc
	UsersAdminsCreate      echo.HandlerFunc
	UsersAdminsUnarchive   echo.HandlerFunc
	UsersAdminsUpdate      echo.HandlerFunc
	UsersTeachersArchive   echo.HandlerFunc
	UsersTeachersCreate    echo.HandlerFunc
	UsersTeachersGet       echo.HandlerFunc
	UsersTeachersList      echo.HandlerFunc
	UsersTeachersUnarchive echo.HandlerFunc
	UsersTeachersUpdate    echo.HandlerFunc

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
	SubjectsGet       echo.HandlerFunc
	SubjectsUpdate    echo.HandlerFunc

	StudentsCreate    echo.HandlerFunc
	StudentsArchive   echo.HandlerFunc
	StudentsUnarchive echo.HandlerFunc
	StudentsList      echo.HandlerFunc
	StudentsGet       echo.HandlerFunc
	StudentsUpdate    echo.HandlerFunc

	GroupsCreate       echo.HandlerFunc
	GroupsList         echo.HandlerFunc
	GroupsGet          echo.HandlerFunc
	GroupsUpdate       echo.HandlerFunc
	GroupsSyncStudents echo.HandlerFunc
	GroupsArchive      echo.HandlerFunc
	GroupsUnarchive    echo.HandlerFunc

	PlansCreate    echo.HandlerFunc
	PlansList      echo.HandlerFunc
	PlansArchive   echo.HandlerFunc
	PlansGet       echo.HandlerFunc
	PlansUnarchive echo.HandlerFunc

	SubscriptionsCreate echo.HandlerFunc
	SubscriptionsList   echo.HandlerFunc

	LessonsIndividualCreate    echo.HandlerFunc
	LessonsGroupCreate         echo.HandlerFunc
	LessonsTemplatesCreate     echo.HandlerFunc
	LessonsTemplatesDeactivate echo.HandlerFunc
	LessonsGet                 echo.HandlerFunc
	LessonsList                echo.HandlerFunc
	LessonsUpdate              echo.HandlerFunc
	LessonsCancel              echo.HandlerFunc
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
		protected.Use(mw.JWT(cfg.JWT.Secret))

		protectedAuthGroup := protected.Group("/auth")
		protectedAuthGroup.GET("/me", h.AuthMe)
		protectedAuthGroup.POST("/logout", h.AuthLogout)

		usersGroup := protected.Group("/users")

		adminsGroup := usersGroup.Group("/admins")
		adminsGroup.Use(mw.RequireRoles(domain.RoleSuperadmin))
		adminsGroup.POST("", h.UsersAdminsCreate)
		adminsGroup.GET("", h.UsersAdminsList)
		adminsGroup.GET("/:id", h.UsersAdminsGet)
		adminsGroup.PUT("/:id", h.UsersAdminsUpdate)
		adminsGroup.PATCH("/:id/archive", h.UsersAdminsArchive)
		adminsGroup.PATCH("/:id/unarchive", h.UsersAdminsUnarchive)

		teachersGroup := usersGroup.Group("/teachers")
		teachersGroup.Use(mw.RequireRoles(domain.RoleSuperadmin, domain.RoleAdmin))
		teachersGroup.GET("", h.UsersTeachersList)
		teachersGroup.GET("/:id", h.UsersTeachersGet)
		teachersGroup.POST("", h.UsersTeachersCreate)
		teachersGroup.PATCH("/:id/archive", h.UsersTeachersArchive)
		teachersGroup.PATCH("/:id/unarchive", h.UsersTeachersUnarchive)
		teachersGroup.PUT("/:id", h.UsersTeachersUpdate)

		branchesGroup := protected.Group("/branches")
		branchesGroup.Use(mw.RequireRoles(domain.RoleSuperadmin, domain.RoleAdmin))
		branchesGroup.POST("", h.BranchesCreate)
		branchesGroup.GET("", h.BranchesList)
		branchesGroup.GET("/:id", h.BranchesGet)
		branchesGroup.PUT("/:id", h.BranchesUpdate)
		branchesGroup.PATCH("/:id/archive", h.BranchesArchive)
		branchesGroup.PATCH("/:id/unarchive", h.BranchesUnarchive)

		subjectsGroup := protected.Group("/subjects")
		subjectsGroup.Use(mw.RequireRoles(domain.RoleSuperadmin, domain.RoleAdmin))
		subjectsGroup.POST("", h.SubjectsCreate)
		subjectsGroup.GET("", h.SubjectsList)
		subjectsGroup.GET("/:id", h.SubjectsGet)
		subjectsGroup.PUT("/:id", h.SubjectsUpdate)
		subjectsGroup.PATCH("/:id/archive", h.SubjectsArchive)
		subjectsGroup.PATCH("/:id/unarchive", h.SubjectsUnarchive)

		studentsGroup := protected.Group("/students")
		studentsGroup.POST("", h.StudentsCreate, mw.RequireRoles(domain.RoleSuperadmin, domain.RoleAdmin))
		studentsGroup.GET("", h.StudentsList, mw.RequireRoles(domain.RoleSuperadmin, domain.RoleAdmin, domain.RoleTeacher))
		studentsGroup.GET("/:id", h.StudentsGet, mw.RequireRoles(domain.RoleSuperadmin, domain.RoleAdmin, domain.RoleTeacher))
		studentsGroup.PUT("/:id", h.StudentsUpdate, mw.RequireRoles(domain.RoleSuperadmin, domain.RoleAdmin))
		studentsGroup.PATCH("/:id/archive", h.StudentsArchive, mw.RequireRoles(domain.RoleSuperadmin, domain.RoleAdmin))
		studentsGroup.PATCH("/:id/unarchive", h.StudentsUnarchive, mw.RequireRoles(domain.RoleSuperadmin, domain.RoleAdmin))
		studentsGroup.POST("/:id/subscriptions", h.SubscriptionsCreate, mw.RequireRoles(domain.RoleSuperadmin, domain.RoleAdmin))
		studentsGroup.GET("/:id/subscriptions", h.SubscriptionsList, mw.RequireRoles(domain.RoleSuperadmin, domain.RoleAdmin))

		groupsGroup := protected.Group("/groups")
		groupsGroup.POST("", h.GroupsCreate, mw.RequireRoles(domain.RoleSuperadmin, domain.RoleAdmin))
		groupsGroup.GET("", h.GroupsList, mw.RequireRoles(domain.RoleSuperadmin, domain.RoleAdmin, domain.RoleTeacher))
		groupsGroup.GET("/:id", h.GroupsGet, mw.RequireRoles(domain.RoleSuperadmin, domain.RoleAdmin, domain.RoleTeacher))
		groupsGroup.PUT("/:id", h.GroupsUpdate, mw.RequireRoles(domain.RoleSuperadmin, domain.RoleAdmin))
		groupsGroup.PUT("/:id/students", h.GroupsSyncStudents, mw.RequireRoles(domain.RoleSuperadmin, domain.RoleAdmin))
		groupsGroup.PATCH("/:id/archive", h.GroupsArchive, mw.RequireRoles(domain.RoleSuperadmin, domain.RoleAdmin))
		groupsGroup.PATCH("/:id/unarchive", h.GroupsUnarchive, mw.RequireRoles(domain.RoleSuperadmin, domain.RoleAdmin))

		plansGroup := protected.Group("/plans")
		plansGroup.Use(mw.RequireRoles(domain.RoleSuperadmin, domain.RoleAdmin))
		plansGroup.POST("", h.PlansCreate)
		plansGroup.GET("", h.PlansList)
		plansGroup.GET("/:id", h.PlansGet)
		plansGroup.PATCH("/:id/archive", h.PlansArchive)
		plansGroup.PATCH("/:id/unarchive", h.PlansUnarchive)

		lessonsGroup := protected.Group("/lessons")
		lessonsGroup.Use(mw.RequireRoles(domain.RoleSuperadmin, domain.RoleAdmin))
		lessonsGroup.POST("/individual", h.LessonsIndividualCreate)
		lessonsGroup.POST("/group", h.LessonsGroupCreate)
		lessonsGroup.POST("/templates", h.LessonsTemplatesCreate)
		lessonsGroup.PATCH("/templates/:id/deactivate", h.LessonsTemplatesDeactivate)
		lessonsGroup.PATCH("/:id", h.LessonsUpdate)
		lessonsGroup.PATCH("/:id/cancel", h.LessonsCancel)

		lessonsReadGroup := protected.Group("/lessons")
		lessonsReadGroup.Use(mw.RequireRoles(domain.RoleSuperadmin, domain.RoleAdmin, domain.RoleTeacher))

		lessonsReadGroup.GET("", h.LessonsList)
		lessonsReadGroup.GET("/:id", h.LessonsGet)
	}
}
