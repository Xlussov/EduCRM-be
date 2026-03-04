package http

import (
	"net/http"
	"time"

	mw "github.com/Xlussov/EduCRM-be/internal/controller/http/middleware"
	"github.com/Xlussov/EduCRM-be/pkg/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Error(msg string)
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Errorf(format string, args ...any)
}

func Init(log Logger, cfg *config.Config, e *echo.Echo) {
	e.Use(mw.RequestLogger(log))
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.CORS())

	e.GET("/ping", func(c echo.Context) error {
		time.Sleep(1 * time.Second)
		return c.String(http.StatusOK, "pong")
	})

	//todo
	// api := e.Group("/api")
	// v1 := api.Group("/v1")

}
