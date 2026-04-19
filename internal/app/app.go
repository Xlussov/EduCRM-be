package app

import (
	"context"
	"net/http"
	"sync"

	"github.com/Xlussov/EduCRM-be/internal/adapter/postgres/postgres"
	repo "github.com/Xlussov/EduCRM-be/internal/adapter/postgres/repos"
	"github.com/Xlussov/EduCRM-be/internal/auth/login"
	logout "github.com/Xlussov/EduCRM-be/internal/auth/logut"
	"github.com/Xlussov/EduCRM-be/internal/auth/me"
	"github.com/Xlussov/EduCRM-be/internal/auth/refresh"
	branchesarchive "github.com/Xlussov/EduCRM-be/internal/branches/archive"
	branchescreate "github.com/Xlussov/EduCRM-be/internal/branches/create"
	branchesget "github.com/Xlussov/EduCRM-be/internal/branches/get"
	brancheslist "github.com/Xlussov/EduCRM-be/internal/branches/list"
	branchesunarchive "github.com/Xlussov/EduCRM-be/internal/branches/unarchive"
	branchesupdate "github.com/Xlussov/EduCRM-be/internal/branches/update"
	httprouter "github.com/Xlussov/EduCRM-be/internal/controller/http"
	groupsarchive "github.com/Xlussov/EduCRM-be/internal/groups/archive"
	groupscreate "github.com/Xlussov/EduCRM-be/internal/groups/create"
	groupsget "github.com/Xlussov/EduCRM-be/internal/groups/get"
	groupslist "github.com/Xlussov/EduCRM-be/internal/groups/list"
	groupssyncstudents "github.com/Xlussov/EduCRM-be/internal/groups/sync_students"
	groupsunarchive "github.com/Xlussov/EduCRM-be/internal/groups/unarchive"
	groupsupdate "github.com/Xlussov/EduCRM-be/internal/groups/update"
	lessonscancel "github.com/Xlussov/EduCRM-be/internal/lessons/cancel"
	lessonsgroupcreate "github.com/Xlussov/EduCRM-be/internal/lessons/create_group"
	lessonsindividualcreate "github.com/Xlussov/EduCRM-be/internal/lessons/create_individual"
	lessonstemplatescreate "github.com/Xlussov/EduCRM-be/internal/lessons/create_template"
	lessonslist "github.com/Xlussov/EduCRM-be/internal/lessons/list"
	lessonsupdate "github.com/Xlussov/EduCRM-be/internal/lessons/update"
	plansarchive "github.com/Xlussov/EduCRM-be/internal/plans/archive"
	planscreate "github.com/Xlussov/EduCRM-be/internal/plans/create"
	plansget "github.com/Xlussov/EduCRM-be/internal/plans/get"
	planslist "github.com/Xlussov/EduCRM-be/internal/plans/list"
	plansunarchive "github.com/Xlussov/EduCRM-be/internal/plans/unarchive"
	studentsarchive "github.com/Xlussov/EduCRM-be/internal/students/archive"
	studentscreate "github.com/Xlussov/EduCRM-be/internal/students/create"
	studentsget "github.com/Xlussov/EduCRM-be/internal/students/get"
	studentslist "github.com/Xlussov/EduCRM-be/internal/students/list"
	studentsunarchive "github.com/Xlussov/EduCRM-be/internal/students/unarchive"
	studentsupdate "github.com/Xlussov/EduCRM-be/internal/students/update"
	subjectsarchive "github.com/Xlussov/EduCRM-be/internal/subjects/archive"
	subjectscreate "github.com/Xlussov/EduCRM-be/internal/subjects/create"
	subjectsget "github.com/Xlussov/EduCRM-be/internal/subjects/get"
	subjectslist "github.com/Xlussov/EduCRM-be/internal/subjects/list"
	subjectsunarchive "github.com/Xlussov/EduCRM-be/internal/subjects/unarchive"
	subjectsupdate "github.com/Xlussov/EduCRM-be/internal/subjects/update"
	subscriptionscreate "github.com/Xlussov/EduCRM-be/internal/subscriptions/create"
	subscriptionslist "github.com/Xlussov/EduCRM-be/internal/subscriptions/list"
	adminarchive "github.com/Xlussov/EduCRM-be/internal/users/admins/archive"
	admincreate "github.com/Xlussov/EduCRM-be/internal/users/admins/create"
	adminget "github.com/Xlussov/EduCRM-be/internal/users/admins/get"
	adminlist "github.com/Xlussov/EduCRM-be/internal/users/admins/list"
	adminunarchive "github.com/Xlussov/EduCRM-be/internal/users/admins/unarchive"
	adminupdate "github.com/Xlussov/EduCRM-be/internal/users/admins/update"
	teacherarchive "github.com/Xlussov/EduCRM-be/internal/users/teachers/archive"
	teachercreate "github.com/Xlussov/EduCRM-be/internal/users/teachers/create"
	teacherget "github.com/Xlussov/EduCRM-be/internal/users/teachers/get"
	teacherlist "github.com/Xlussov/EduCRM-be/internal/users/teachers/list"
	teacherunarchive "github.com/Xlussov/EduCRM-be/internal/users/teachers/unarchive"
	teacherupdate "github.com/Xlussov/EduCRM-be/internal/users/teachers/update"
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

	dbPool, err := postgres.New(ctx, cfg.Postgres.URL, log)
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
	scheduleRepo := repo.NewScheduleRepository(dbPool.Conn())

	txManager := postgres.NewTxManager(dbPool.Conn())

	loginUC := login.NewUseCase(userRepo, authRepo, cfg.JWT.Secret, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL)
	refreshUC := refresh.NewUseCase(userRepo, authRepo, cfg.JWT.Secret, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL)
	logoutUC := logout.NewUseCase(authRepo)
	adminsArchiveUC := adminarchive.NewUseCase(userRepo)
	adminsCreateUC := admincreate.NewUseCase(userRepo, txManager)
	adminsGetUC := adminget.NewUseCase(userRepo)
	adminsListUC := adminlist.NewUseCase(userRepo)
	adminsUnarchiveUC := adminunarchive.NewUseCase(userRepo)
	adminsUpdateUC := adminupdate.NewUseCase(userRepo, txManager)
	teachersArchiveUC := teacherarchive.NewUseCase(userRepo)
	teachersCreateUC := teachercreate.NewUseCase(userRepo, txManager)
	teachersGetUC := teacherget.NewUseCase(userRepo)
	teachersListUC := teacherlist.NewUseCase(userRepo)
	teachersUnarchiveUC := teacherunarchive.NewUseCase(userRepo)
	teachersUpdateUC := teacherupdate.NewUseCase(userRepo, scheduleRepo, txManager)
	meUC := me.NewUseCase(userRepo)

	branchesCreateUC := branchescreate.NewUseCase(branchRepo, userRepo, txManager)
	branchesArchiveUC := branchesarchive.NewUseCase(branchRepo)
	branchesListUC := brancheslist.NewUseCase(branchRepo)
	branchesGetUC := branchesget.NewUseCase(branchRepo)
	branchesUpdateUC := branchesupdate.NewUseCase(branchRepo)
	branchesUnarchiveUC := branchesunarchive.NewUseCase(branchRepo)

	subjectsCreateUC := subjectscreate.NewUseCase(subjectRepo, branchRepo)
	subjectsArchiveUC := subjectsarchive.NewUseCase(subjectRepo)
	subjectsListUC := subjectslist.NewUseCase(subjectRepo)
	subjectsGetUC := subjectsget.NewUseCase(subjectRepo)
	subjectsUpdateUC := subjectsupdate.NewUseCase(subjectRepo)
	subjectsUnarchiveUC := subjectsunarchive.NewUseCase(subjectRepo)

	studentsCreateUC := studentscreate.NewUseCase(studentRepo, userRepo)
	studentsArchiveUC := studentsarchive.NewUseCase(studentRepo)
	studentUnarchiveUC := studentsunarchive.NewUseCase(studentRepo)
	studentsListUC := studentslist.NewUseCase(studentRepo)
	studentsGetUC := studentsget.NewUseCase(studentRepo)
	studentsUpdateUC := studentsupdate.NewUseCase(studentRepo)

	groupsCreateUC := groupscreate.NewUseCase(groupRepo, userRepo)
	groupsListUC := groupslist.NewUseCase(groupRepo)
	groupsGetUC := groupsget.NewUseCase(groupRepo)
	groupsUpdateUC := groupsupdate.NewUseCase(groupRepo)
	groupsSyncStudentsUC := groupssyncstudents.NewUseCase(groupRepo, studentRepo, txManager)
	groupsArchiveUC := groupsarchive.NewUseCase(groupRepo)
	groupsUnarchiveUC := groupsunarchive.NewUseCase(groupRepo)

	plansCreateUC := planscreate.NewUseCase(txManager, planRepo, userRepo)
	plansListUC := planslist.NewUseCase(planRepo)
	plansArchiveUC := plansarchive.NewUseCase(planRepo)
	plansGetUC := plansget.NewUseCase(planRepo)
	plansUnarchiveUC := plansunarchive.NewUseCase(planRepo)

	subscriptionsCreateUC := subscriptionscreate.NewUseCase(planRepo, studentRepo)
	subscriptionsListUC := subscriptionslist.NewUseCase(planRepo, studentRepo)

	lessonsIndividualCreateUC := lessonsindividualcreate.NewUseCase(scheduleRepo, userRepo)
	lessonsGroupCreateUC := lessonsgroupcreate.NewUseCase(scheduleRepo, groupRepo, userRepo)
	lessonsTemplatesCreateUC := lessonstemplatescreate.NewUseCase(scheduleRepo, groupRepo, userRepo)
	lessonsListUC := lessonslist.NewUseCase(scheduleRepo)
	lessonsUpdateUC := lessonsupdate.NewUseCase(scheduleRepo, groupRepo, userRepo)
	lessonsCancelUC := lessonscancel.NewUseCase(scheduleRepo)

	h := httprouter.Handlers{
		AuthLogin:              login.NewHandler(loginUC).Handle,
		AuthRefresh:            refresh.NewHandler(refreshUC).Handle,
		AuthLogout:             logout.NewHandler(logoutUC).Handle,
		UsersAdminsArchive:     adminarchive.NewHandler(adminsArchiveUC).Handle,
		UsersAdminsGet:         adminget.NewHandler(adminsGetUC).Handle,
		UsersAdminsList:        adminlist.NewHandler(adminsListUC).Handle,
		UsersAdminsCreate:      admincreate.NewHandler(adminsCreateUC).Handle,
		UsersAdminsUnarchive:   adminunarchive.NewHandler(adminsUnarchiveUC).Handle,
		UsersAdminsUpdate:      adminupdate.NewHandler(adminsUpdateUC).Handle,
		UsersTeachersArchive:   teacherarchive.NewHandler(teachersArchiveUC).Handle,
		UsersTeachersCreate:    teachercreate.NewHandler(teachersCreateUC).Handle,
		UsersTeachersGet:       teacherget.NewHandler(teachersGetUC).Handle,
		UsersTeachersList:      teacherlist.NewHandler(teachersListUC).Handle,
		UsersTeachersUnarchive: teacherunarchive.NewHandler(teachersUnarchiveUC).Handle,
		UsersTeachersUpdate:    teacherupdate.NewHandler(teachersUpdateUC).Handle,
		AuthMe:                 me.NewHandler(meUC).Handle,

		BranchesCreate:    branchescreate.NewHandler(branchesCreateUC).Handle,
		BranchesArchive:   branchesarchive.NewHandler(branchesArchiveUC).Handle,
		BranchesUnarchive: branchesunarchive.NewHandler(branchesUnarchiveUC).Handle,
		BranchesList:      brancheslist.NewHandler(branchesListUC).Handle,
		BranchesGet:       branchesget.NewHandler(branchesGetUC).Handle,
		BranchesUpdate:    branchesupdate.NewHandler(branchesUpdateUC).Handle,

		SubjectsCreate:    subjectscreate.NewHandler(subjectsCreateUC).Handle,
		SubjectsArchive:   subjectsarchive.NewHandler(subjectsArchiveUC).Handle,
		SubjectsUnarchive: subjectsunarchive.NewHandler(subjectsUnarchiveUC).Handle,
		SubjectsList:      subjectslist.NewHandler(subjectsListUC).Handle,
		SubjectsGet:       subjectsget.NewHandler(subjectsGetUC).Handle,
		SubjectsUpdate:    subjectsupdate.NewHandler(subjectsUpdateUC).Handle,

		StudentsCreate:    studentscreate.NewHandler(studentsCreateUC).Handle,
		StudentsArchive:   studentsarchive.NewHandler(studentsArchiveUC).Handle,
		StudentsUnarchive: studentsunarchive.NewHandler(studentUnarchiveUC).Handle,
		StudentsList:      studentslist.NewHandler(studentsListUC).Handle,
		StudentsGet:       studentsget.NewHandler(studentsGetUC).Handle,
		StudentsUpdate:    studentsupdate.NewHandler(studentsUpdateUC).Handle,

		GroupsCreate:       groupscreate.NewHandler(groupsCreateUC).Handle,
		GroupsList:         groupslist.NewHandler(groupsListUC).Handle,
		GroupsGet:          groupsget.NewHandler(groupsGetUC).Handle,
		GroupsUpdate:       groupsupdate.NewHandler(groupsUpdateUC).Handle,
		GroupsSyncStudents: groupssyncstudents.NewHandler(groupsSyncStudentsUC).Handle,
		GroupsArchive:      groupsarchive.NewHandler(groupsArchiveUC).Handle,
		GroupsUnarchive:    groupsunarchive.NewHandler(groupsUnarchiveUC).Handle,

		PlansCreate:    planscreate.NewHandler(plansCreateUC).Handle,
		PlansList:      planslist.NewHandler(plansListUC).Handle,
		PlansArchive:   plansarchive.NewHandler(plansArchiveUC).Handle,
		PlansGet:       plansget.NewHandler(plansGetUC).Handle,
		PlansUnarchive: plansunarchive.NewHandler(plansUnarchiveUC).Handle,

		SubscriptionsCreate: subscriptionscreate.NewHandler(subscriptionsCreateUC).Handle,
		SubscriptionsList:   subscriptionslist.NewHandler(subscriptionsListUC).Handle,

		LessonsIndividualCreate: lessonsindividualcreate.NewHandler(lessonsIndividualCreateUC).Handle,
		LessonsGroupCreate:      lessonsgroupcreate.NewHandler(lessonsGroupCreateUC).Handle,
		LessonsTemplatesCreate:  lessonstemplatescreate.NewHandler(lessonsTemplatesCreateUC).Handle,
		LessonsList:             lessonslist.NewHandler(lessonsListUC).Handle,
		LessonsUpdate:           lessonsupdate.NewHandler(lessonsUpdateUC).Handle,
		LessonsCancel:           lessonscancel.NewHandler(lessonsCancelUC).Handle,
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
