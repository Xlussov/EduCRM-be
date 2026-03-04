package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
)

type Logger interface {
	Infof(format string, args ...any)
}

func RequestLogger(l Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)

			req := c.Request()
			res := c.Response()

			duration := time.Since(start)
			ms := float64(duration.Microseconds()) / 1000.0

			ip := c.RealIP()

			l.Infof("%d %s %s %.1fms IP:%s",
				res.Status,
				req.Method,
				req.URL.Path,
				ms,
				ip,
			)

			return err
		}
	}
}
