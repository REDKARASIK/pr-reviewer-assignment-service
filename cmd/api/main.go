package main

import (
	"context"
	"log"
	http2 "net/http"
	"os"
	"pr-reviewer-assigment-service/docs"
	app2 "pr-reviewer-assigment-service/internal/app"
	"pr-reviewer-assigment-service/internal/http"
	"time"
)

func main() {
	docs.SwaggerInfo.BasePath = "/"
	dsn := os.Getenv("DB_DSN")
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

	log.Println("Starting server on :8080")
	if err := http2.ListenAndServe(":8080", server); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
