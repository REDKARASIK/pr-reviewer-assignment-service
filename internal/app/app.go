package app

import (
	"context"
	"net/http"
	"pr-reviewer-assigment-service/internal/http/router"
	"pr-reviewer-assigment-service/internal/http/v1/users"
	"pr-reviewer-assigment-service/internal/repository/postgres"
	"pr-reviewer-assigment-service/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	db          *pgxpool.Pool
	Router      *router.Router
	UserHandler *users.UsersHandler
}

func NewApp(ctx context.Context, dsn string) (*App, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	// repo
	userRepo := postgres.NewUserRepository(pool)
	prRepo := postgres.NewPullRequestRepository(pool)

	// service
	userServ := service.NewUserService(userRepo)
	prServ := service.NewPullRequestService(prRepo)

	// handlers
	userHandler := users.NewUsersHandler(userServ, prServ)

	app := &App{
		db:          pool,
		UserHandler: userHandler,
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
