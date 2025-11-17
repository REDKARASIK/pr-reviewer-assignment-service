package main

import (
	"context"
	"fmt"
	"log"
	http2 "net/http"
	"pr-reviewer-assigment-service/docs"
	app2 "pr-reviewer-assigment-service/internal/app"
	"pr-reviewer-assigment-service/internal/config"
	"pr-reviewer-assigment-service/internal/http"
	"time"
)

func main() {
	docs.SwaggerInfo.BasePath = "/"

	cfg, err := config.Load("config.toml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	dsn := cfg.Postgres.DSN()
	if dsn == "" {
		log.Fatal("DB_DSN environment variable not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	app, err := app2.NewApp(ctx, dsn)
	if err != nil {
		log.Fatalf("failed to init app: %v", err)
	}
	defer app.Close()

	server := http.RegisterRoutes(http.RoutesHandlers{
		Router:       app.Router,
		UserHandler:  app.UserHandler,
		TeamHandler:  app.TeamHandler,
		PrHandler:    app.PRHandler,
		StatsHandler: app.StatsHandler,
	})

	addr := fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port)

	log.Println("Starting server on", addr)
	if err := http2.ListenAndServe(addr, server); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
