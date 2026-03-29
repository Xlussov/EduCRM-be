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
		branchesGroup.PATCH("/:id/archive", h.BranchesArchive)
	}
}
