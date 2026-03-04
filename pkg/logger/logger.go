// package logger

// import (
// 	"fmt"
// 	"log/slog"
// 	"os"

// 	"github.com/Xlussov/dmkt-bot-tg/pkg/config"
// )

// const (
// 	envLocal = "local"
// 	envDev   = "dev"
// 	envProd  = "prod"
// )
// type Logger struct {
// 	log *slog.Logger
// }

// func (s *Logger) Debugf(format string, args ...any) {
// 	s.log.Debug(fmt.Sprintf(format, args...))
// }

// func (s *Logger) Infof(format string, args ...any) {
// 	s.log.Info(fmt.Sprintf(format, args...))
// }

// func (s *Logger) Errorf(format string, args ...any) {
// 	s.log.Error(fmt.Sprintf(format, args...))
// }

// func New(cfg *config.Config) *slog.Logger {
// 	var log *slog.Logger

// 	switch cfg.Env {
// 	case envLocal:
// 		log = slog.New(
// 			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
// 		)
// 	case envDev:
// 		log = slog.New(
// 			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
// 		)
// 	case envProd:
// 		log = slog.New(
// 			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
// 		)
// 	}

// 	return log
// }

package logger

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/Xlussov/EduCRM-be/pkg/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

type Logger struct {
	log *slog.Logger
}

func (l *Logger) Debug(massage string) {
	l.log.Debug(massage)
}

func (l *Logger) Info(massage string) {
	l.log.Info(massage)
}

func (l *Logger) Error(massage string) {
	l.log.Error(massage)
}

func (l *Logger) Debugf(format string, args ...any) {
	l.log.Debug(fmt.Sprintf(format, args...))
}

func (l *Logger) Infof(format string, args ...any) {
	l.log.Info(fmt.Sprintf(format, args...))
}

func (l *Logger) Errorf(format string, args ...any) {
	l.log.Error(fmt.Sprintf(format, args...))
}

func New(cfg *config.Config) *Logger {
	var slogLogger *slog.Logger

	switch cfg.Env {
	case envLocal:
		slogLogger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		slogLogger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		slogLogger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		slogLogger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	}

	return &Logger{log: slogLogger}
}
