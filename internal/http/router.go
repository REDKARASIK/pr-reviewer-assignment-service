package http

import (
	"net/http"
	"pr-reviewer-assigment-service/internal/http/router"
)

type RoutesHandlers struct {
	Router *router.Router
}

func RegisterRoutes(h RoutesHandlers) http.Handler {
	r := h.Router

	// users := r.Group("/users")

	return r.Handler()
}
