package app

import (
	"context"
	"net/http"
	"sync"

	"github.com/Xlussov/EduCRM-be/internal/adapter/postgres/postgres"
	repo "github.com/Xlussov/EduCRM-be/internal/adapter/postgres/repos"
	"github.com/Xlussov/EduCRM-be/internal/auth/login"
	"github.com/Xlussov/EduCRM-be/internal/auth/refresh"
	branchesarchive "github.com/Xlussov/EduCRM-be/internal/branches/archive"
	branchescreate "github.com/Xlussov/EduCRM-be/internal/branches/create"
	branchesget "github.com/Xlussov/EduCRM-be/internal/branches/get"
	brancheslist "github.com/Xlussov/EduCRM-be/internal/branches/list"
	branchesupdate "github.com/Xlussov/EduCRM-be/internal/branches/update"
	httprouter "github.com/Xlussov/EduCRM-be/internal/controller/http"
	groupsaddstudents "github.com/Xlussov/EduCRM-be/internal/groups/add_students"
	groupscreate "github.com/Xlussov/EduCRM-be/internal/groups/create"
	groupsget "github.com/Xlussov/EduCRM-be/internal/groups/get"
	groupslist "github.com/Xlussov/EduCRM-be/internal/groups/list"
	groupsremovestudent "github.com/Xlussov/EduCRM-be/internal/groups/remove_student"
	groupsupdate "github.com/Xlussov/EduCRM-be/internal/groups/update"
	plansarchive "github.com/Xlussov/EduCRM-be/internal/plans/archive"
	planscreate "github.com/Xlussov/EduCRM-be/internal/plans/create"
	planslist "github.com/Xlussov/EduCRM-be/internal/plans/list"
	studentsarchive "github.com/Xlussov/EduCRM-be/internal/students/archive"
	studentscreate "github.com/Xlussov/EduCRM-be/internal/students/create"
	studentsget "github.com/Xlussov/EduCRM-be/internal/students/get"
	studentslist "github.com/Xlussov/EduCRM-be/internal/students/list"
	studentsupdate "github.com/Xlussov/EduCRM-be/internal/students/update"
	subjectsarchive "github.com/Xlussov/EduCRM-be/internal/subjects/archive"
	subjectscreate "github.com/Xlussov/EduCRM-be/internal/subjects/create"
	subjectslist "github.com/Xlussov/EduCRM-be/internal/subjects/list"
	subjectsupdate "github.com/Xlussov/EduCRM-be/internal/subjects/update"
	subscriptionscreate "github.com/Xlussov/EduCRM-be/internal/subscriptions/create"
	subscriptionslist "github.com/Xlussov/EduCRM-be/internal/subscriptions/list"
	"github.com/Xlussov/EduCRM-be/internal/users/admins"
	"github.com/Xlussov/EduCRM-be/internal/users/teachers"
	"github.com/Xlussov/EduCRM-be/pkg/config"
	"github.com/Xlussov/EduCRM-be/pkg/validator"
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
	studentRepo := repo.NewStudentRepository(dbPool.Conn())
	groupRepo := repo.NewGroupRepository(dbPool.Conn())
	planRepo := repo.NewSubscriptionRepository(dbPool.Conn())

	txManager := postgres.NewTxManager(dbPool.Conn())

	loginUC := login.NewUseCase(userRepo, authRepo, cfg.JWTSecret)
	refreshUC := refresh.NewUseCase(userRepo, authRepo, cfg.JWTSecret)
	adminsUC := admins.NewUseCase(userRepo, txManager)
	teachersUC := teachers.NewUseCase(userRepo, txManager)
	branchesCreateUC := branchescreate.NewUseCase(branchRepo, userRepo, txManager)
	branchesArchiveUC := branchesarchive.NewUseCase(branchRepo)
	branchesListUC := brancheslist.NewUseCase(branchRepo)
	branchesGetUC := branchesget.NewUseCase(branchRepo)
	branchesUpdateUC := branchesupdate.NewUseCase(branchRepo)
	subjectsCreateUC := subjectscreate.NewUseCase(subjectRepo)
	subjectsArchiveUC := subjectsarchive.NewUseCase(subjectRepo)
	subjectsListUC := subjectslist.NewUseCase(subjectRepo)
	subjectsUpdateUC := subjectsupdate.NewUseCase(subjectRepo)
	studentsCreateUC := studentscreate.NewUseCase(studentRepo, userRepo)
	studentsArchiveUC := studentsarchive.NewUseCase(studentRepo)
	studentsListUC := studentslist.NewUseCase(studentRepo, userRepo)
	studentsGetUC := studentsget.NewUseCase(studentRepo)
	studentsUpdateUC := studentsupdate.NewUseCase(studentRepo, userRepo)
	groupsCreateUC := groupscreate.NewUseCase(groupRepo, userRepo)
	groupsListUC := groupslist.NewUseCase(groupRepo, userRepo)
	groupsGetUC := groupsget.NewUseCase(groupRepo, userRepo)
	groupsUpdateUC := groupsupdate.NewUseCase(groupRepo, userRepo)
	groupsAddStudentsUC := groupsaddstudents.NewUseCase(groupRepo, userRepo, studentRepo, txManager)
	groupsRemoveStudentUC := groupsremovestudent.NewUseCase(groupRepo, userRepo, studentRepo)
	plansCreateUC := planscreate.NewUseCase(txManager, planRepo, userRepo)
	plansListUC := planslist.NewUseCase(planRepo, userRepo)
	plansArchiveUC := plansarchive.NewUseCase(planRepo, userRepo)
	subscriptionsCreateUC := subscriptionscreate.NewUseCase(planRepo, userRepo, studentRepo)
	subscriptionsListUC := subscriptionslist.NewUseCase(planRepo, userRepo, studentRepo)

	h := httprouter.Handlers{
		AuthLogin:           login.NewHandler(loginUC).Handle,
		AuthRefresh:         refresh.NewHandler(refreshUC).Handle,
		UsersAdmins:         admins.NewHandler(adminsUC).Handle,
		UsersTeachers:       teachers.NewHandler(teachersUC).Handle,
		BranchesCreate:      branchescreate.NewHandler(branchesCreateUC).Handle,
		BranchesArchive:     branchesarchive.NewHandler(branchesArchiveUC).Handle,
		BranchesList:        brancheslist.NewHandler(branchesListUC).Handle,
		BranchesGet:         branchesget.NewHandler(branchesGetUC).Handle,
		BranchesUpdate:      branchesupdate.NewHandler(branchesUpdateUC).Handle,
		SubjectsCreate:      subjectscreate.NewHandler(subjectsCreateUC).Handle,
		SubjectsArchive:     subjectsarchive.NewHandler(subjectsArchiveUC).Handle,
		SubjectsList:        subjectslist.NewHandler(subjectsListUC).Handle,
		SubjectsUpdate:      subjectsupdate.NewHandler(subjectsUpdateUC).Handle,
		StudentsCreate:      studentscreate.NewHandler(studentsCreateUC).Handle,
		StudentsArchive:     studentsarchive.NewHandler(studentsArchiveUC).Handle,
		StudentsList:        studentslist.NewHandler(studentsListUC).Handle,
		StudentsGet:         studentsget.NewHandler(studentsGetUC).Handle,
		StudentsUpdate:      studentsupdate.NewHandler(studentsUpdateUC).Handle,
		GroupsCreate:        groupscreate.NewHandler(groupsCreateUC).Handle,
		GroupsList:          groupslist.NewHandler(groupsListUC).Handle,
		GroupsGet:           groupsget.NewHandler(groupsGetUC).Handle,
		GroupsUpdate:        groupsupdate.NewHandler(groupsUpdateUC).Handle,
		GroupsAddStudents:   groupsaddstudents.NewHandler(groupsAddStudentsUC).Handle,
		GroupsRemoveStudent: groupsremovestudent.NewHandler(groupsRemoveStudentUC).Handle,
		PlansCreate:         planscreate.NewHandler(plansCreateUC).Handle,
		PlansList:           planslist.NewHandler(plansListUC).Handle,
		SubscriptionsCreate: subscriptionscreate.NewHandler(subscriptionsCreateUC).Handle,
		SubscriptionsList:   subscriptionslist.NewHandler(subscriptionsListUC).Handle,
		PlansArchive:        plansarchive.NewHandler(plansArchiveUC).Handle,
	}

	e.Validator = validator.New()
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
