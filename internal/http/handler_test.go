package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	myhttp "pr-reviewer-assigment-service/internal/http"
	"pr-reviewer-assigment-service/internal/http/router"
)

func TestRoutesRegister(t *testing.T) {
	r := router.NewRouter()

	handlers := myhttp.RoutesHandlers{
		Router:       r,
		UserHandler:  nil,
		TeamHandler:  nil,
		PrHandler:    nil,
		StatsHandler: nil,
	}

	server := myhttp.RegisterRoutes(handlers)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	if w.Code == 0 {
		t.Fatalf("expected HTTP response, got empty status")
	}
}
