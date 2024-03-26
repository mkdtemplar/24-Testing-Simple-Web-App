package main

import (
	"net/http"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func Test_application_routes(t *testing.T) {
	tests := []struct {
		route  string
		method string
	}{
		{route: "/", method: "GET"},
		{route: "/static/*", method: "GET"},
		{route: "/login", method: "POST"},
		{route: "/user/profile", method: "GET"},
	}

	mux := app.routes()

	chiRoutes := mux.(chi.Routes)
	for _, r := range tests {
		if !routeExists(r.route, r.method, chiRoutes) {
			t.Errorf("Route %s is not registered", r.route)
		}
	}
}

func routeExists(testRoute string, testMethod string, chiRoutes chi.Routes) bool {

	found := false

	_ = chi.Walk(chiRoutes, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if strings.EqualFold(method, testMethod) && strings.EqualFold(route, testRoute) {
			found = true
		}
		return nil
	})

	return found
}
