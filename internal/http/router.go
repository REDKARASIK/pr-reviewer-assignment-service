package http

import (
	"net/http"
	"pr-reviewer-assigment-service/internal/http/router"
	"pr-reviewer-assigment-service/internal/http/v1/pull_requests"
	"pr-reviewer-assigment-service/internal/http/v1/teams"
	"pr-reviewer-assigment-service/internal/http/v1/users"

	httpSwagger "github.com/swaggo/http-swagger"
)

type RoutesHandlers struct {
	Router      *router.Router
	UserHandler *users.UsersHandler
	TeamHandler *teams.TeamsHandler
	PrHandler   *pull_requests.PullRequestHandler
}

func RegisterRoutes(h RoutesHandlers) http.Handler {
	r := h.Router

	// users
	usersGroup := r.Group("/users")
	usersGroup.GET("/getReview", h.UserHandler.GetReview)
	usersGroup.POST("/setIsActive", h.UserHandler.SetIsActive)

	// teams
	teamsGroup := r.Group("/team")
	teamsGroup.POST("/add", h.TeamHandler.Add)
	teamsGroup.GET("/get", h.TeamHandler.Get)

	// prs
	prGroup := r.Group("/pullRequest")
	prGroup.POST("/create", h.PrHandler.Create)
	prGroup.POST("/merge", h.PrHandler.Merge)

	// swagger
	r.GET("/swagger", httpSwagger.WrapHandler)
	r.GET("/swagger/", httpSwagger.WrapHandler)
	r.GET("/swagger/index.html", httpSwagger.WrapHandler)
	r.GET("/swagger/doc.json", httpSwagger.WrapHandler)

	return r.Handler()
}
