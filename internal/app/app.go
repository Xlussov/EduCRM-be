package app

import (
	"context"
	"net/http"
	"sync"

	"github.com/Xlussov/EduCRM-be/internal/adapter/postgres/postgres"
	httprouter "github.com/Xlussov/EduCRM-be/internal/controller/http"
	"github.com/Xlussov/EduCRM-be/pkg/config"
	"github.com/labstack/echo/v4"
)

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Error(msg string)
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Errorf(format string, args ...any)
}

type App struct {
	cfg  *config.Config
	log  Logger
	wg   sync.WaitGroup
	echo *echo.Echo
	db   *postgres.Pool
}

func New(ctx context.Context, cfg *config.Config, log Logger) (*App, error) {
	e := echo.New()

	httprouter.Init(log, cfg, e)

	pgCfg := postgres.Config{
		User:     cfg.Postgres.User,
		Password: cfg.Postgres.Password,
		Host:     cfg.Postgres.Host,
		Port:     cfg.Postgres.Port,
		DBName:   cfg.Postgres.DBName,
		SSLMode:  cfg.Postgres.SSLMode,
	}

	dbPool, err := postgres.New(ctx, pgCfg, log)
	if err != nil {
		log.Errorf("failed to init postgres: %v", err)
		return nil, err
	}
	log.Info("successfully connected to postgres")

	return &App{
		cfg:  cfg,
		log:  log,
		echo: e,
		db:   dbPool,
	}, nil
}

func (a *App) Start(ctx context.Context) {
	a.log.Info("starting app services...")

	a.wg.Go(func() {
		err := a.echo.Start(a.cfg.HTTPServer.Address)
		if err != nil && err != http.ErrServerClosed {
			a.log.Errorf("http server error: %v", err)
		}
	})

	a.wg.Go(func() {
		<-ctx.Done()
		a.log.Info("context canceled, shutting down components...")
	})
}

func (a *App) Stop(ctx context.Context) error {
	a.log.Info("graceful shutdown started")

	if err := a.echo.Shutdown(ctx); err != nil {
		return err
	}

	if err := a.db; err != nil {
		a.db.Close()
		a.log.Info("postgres pool closed")
	}

	a.wg.Wait()

	return nil
}
