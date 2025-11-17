package app

import (
	"context"
	"net/http"
	"pr-reviewer-assigment-service/internal/http/router"
	"pr-reviewer-assigment-service/internal/http/v1/pull_requests"
	"pr-reviewer-assigment-service/internal/http/v1/statistics"
	"pr-reviewer-assigment-service/internal/http/v1/teams"
	"pr-reviewer-assigment-service/internal/http/v1/users"
	"pr-reviewer-assigment-service/internal/repository/postgres"
	"pr-reviewer-assigment-service/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	db           *pgxpool.Pool
	Router       *router.Router
	UserHandler  *users.UsersHandler
	TeamHandler  *teams.TeamsHandler
	PRHandler    *pull_requests.PullRequestHandler
	StatsHandler *statistics.StatisticsHandler
}

func NewApp(ctx context.Context, dsn string) (*App, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	// repo
	userRepo := postgres.NewUserRepository(pool)
	prRepo := postgres.NewPullRequestRepository(pool)
	teamRepo := postgres.NewTeamRepository(pool)
	statsRepo := postgres.NewStatisticsPostgresRepository(pool)

	// service
	userServ := service.NewUserService(userRepo)
	prServ := service.NewPullRequestService(prRepo, userRepo, teamRepo)
	teamServ := service.NewTeamService(teamRepo, userRepo)
	statsServ := service.NewStatisticsService(statsRepo)

	// handlers
	userHandler := users.NewUsersHandler(userServ, prServ)
	teamHandler := teams.NewTeamsHandler(teamServ)
	prHandler := pull_requests.NewPullRequestHandler(prServ)
	statsHandler := statistics.NewStatisticsHandler(statsServ)

	app := &App{
		db:           pool,
		UserHandler:  userHandler,
		TeamHandler:  teamHandler,
		PRHandler:    prHandler,
		StatsHandler: statsHandler,
	}

	app.Router = router.NewRouter()

	return app, nil
}

func (a *App) Handler() http.Handler {
	return a.Router.Handler()
}

func (a *App) Run(addr string) error {
	return http.ListenAndServe(addr, a.Handler())
}

func (a *App) Close() {
	if a.db != nil {
		a.db.Close()
	}
}
