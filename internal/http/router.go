package http

import (
	"net/http"
	"pr-reviewer-assigment-service/internal/http/router"
	"pr-reviewer-assigment-service/internal/http/v1/teams"
	"pr-reviewer-assigment-service/internal/http/v1/users"

	httpSwagger "github.com/swaggo/http-swagger"
)

type RoutesHandlers struct {
	Router      *router.Router
	UserHandler *users.UsersHandler
	TeamHandler *teams.TeamsHandler
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

	// swagger
	r.GET("/swagger", httpSwagger.WrapHandler)
	r.GET("/swagger/", httpSwagger.WrapHandler)
	r.GET("/swagger/index.html", httpSwagger.WrapHandler)
	r.GET("/swagger/doc.json", httpSwagger.WrapHandler)

	return r.Handler()
}
