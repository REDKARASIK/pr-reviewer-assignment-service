package http

import (
	"net/http"
	"pr-reviewer-assigment-service/internal/http/router"
	"pr-reviewer-assigment-service/internal/http/v1/users"
)

type RoutesHandlers struct {
	Router      *router.Router
	UserHandler *users.UsersHandler
}

func RegisterRoutes(h RoutesHandlers) http.Handler {
	r := h.Router

	usersGroup := r.Group("/users")
	usersGroup.GET("/getReview", h.UserHandler.GetReview)
	usersGroup.POST("/setIsActive", h.UserHandler.SetIsActive)

	return r.Handler()
}
