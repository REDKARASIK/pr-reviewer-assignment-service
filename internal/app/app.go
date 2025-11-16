package app

import (
	"context"
	"net/http"
	"pr-reviewer-assigment-service/internal/http/router"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	db     *pgxpool.Pool
	Router *router.Router
}

func NewApp(ctx context.Context, dsn string) (*App, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	app := &App{
		db: pool,
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
