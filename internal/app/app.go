package app

import (
	"context"
	"net/http"
	"sync"

	repo "github.com/Xlussov/EduCRM-be/internal/adapter/postgres"
	"github.com/Xlussov/EduCRM-be/internal/adapter/postgres/postgres"
	"github.com/Xlussov/EduCRM-be/internal/auth/login"
	"github.com/Xlussov/EduCRM-be/internal/auth/refresh"
	branchesarchive "github.com/Xlussov/EduCRM-be/internal/branches/archive"
	branchescreate "github.com/Xlussov/EduCRM-be/internal/branches/create"
	branchesget "github.com/Xlussov/EduCRM-be/internal/branches/get"
	brancheslist "github.com/Xlussov/EduCRM-be/internal/branches/list"
	branchesupdate "github.com/Xlussov/EduCRM-be/internal/branches/update"
	httprouter "github.com/Xlussov/EduCRM-be/internal/controller/http"
	subjectsarchive "github.com/Xlussov/EduCRM-be/internal/subjects/archive"
	subjectscreate "github.com/Xlussov/EduCRM-be/internal/subjects/create"
	subjectslist "github.com/Xlussov/EduCRM-be/internal/subjects/list"
	subjectsupdate "github.com/Xlussov/EduCRM-be/internal/subjects/update"
	"github.com/Xlussov/EduCRM-be/internal/users/admins"
	"github.com/Xlussov/EduCRM-be/internal/users/teachers"
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

	userRepo := repo.NewUserRepository(dbPool.Conn())
	authRepo := repo.NewAuthRepository(dbPool.Conn())
	branchRepo := repo.NewBranchRepository(dbPool.Conn())
	subjectRepo := repo.NewSubjectRepository(dbPool.Conn())

	loginUC := login.NewUseCase(userRepo, authRepo, cfg.JWTSecret)
	refreshUC := refresh.NewUseCase(userRepo, authRepo, cfg.JWTSecret)
	adminsUC := admins.NewUseCase(userRepo)
	teachersUC := teachers.NewUseCase(userRepo)
	branchesCreateUC := branchescreate.NewUseCase(branchRepo, userRepo)
	branchesArchiveUC := branchesarchive.NewUseCase(branchRepo)
	branchesListUC := brancheslist.NewUseCase(branchRepo)
	branchesGetUC := branchesget.NewUseCase(branchRepo)
	branchesUpdateUC := branchesupdate.NewUseCase(branchRepo)
	subjectsCreateUC := subjectscreate.NewUseCase(subjectRepo)
	subjectsArchiveUC := subjectsarchive.NewUseCase(subjectRepo)
	subjectsListUC := subjectslist.NewUseCase(subjectRepo)
	subjectsUpdateUC := subjectsupdate.NewUseCase(subjectRepo)

	h := httprouter.Handlers{
		AuthLogin:       login.NewHandler(loginUC).Handle,
		AuthRefresh:     refresh.NewHandler(refreshUC).Handle,
		UsersAdmins:     admins.NewHandler(adminsUC).Handle,
		UsersTeachers:   teachers.NewHandler(teachersUC).Handle,
		BranchesCreate:  branchescreate.NewHandler(branchesCreateUC).Handle,
		BranchesArchive: branchesarchive.NewHandler(branchesArchiveUC).Handle,
		BranchesList:    brancheslist.NewHandler(branchesListUC).Handle,
		BranchesGet:     branchesget.NewHandler(branchesGetUC).Handle,
		BranchesUpdate:  branchesupdate.NewHandler(branchesUpdateUC).Handle,
		SubjectsCreate:  subjectscreate.NewHandler(subjectsCreateUC).Handle,
		SubjectsArchive: subjectsarchive.NewHandler(subjectsArchiveUC).Handle,
		SubjectsList:    subjectslist.NewHandler(subjectsListUC).Handle,
		SubjectsUpdate:  subjectsupdate.NewHandler(subjectsUpdateUC).Handle,
	}

	httprouter.Init(log, cfg, e, h)

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
