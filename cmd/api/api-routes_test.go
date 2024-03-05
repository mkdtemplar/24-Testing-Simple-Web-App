package main

import (
	"testing"
)

func Test_application_routes(t *testing.T) {
	tests := []struct {
		route  string
		method string
	}{
		{route: "/auth", method: "POST"},
		{route: "/refresh-token", method: "POST"},
	}

	mux := app.routes()
	exists := make(map[string]bool)

	for _, routeInfo := range mux.Routes() {
		key := routeInfo.Path + " " + routeInfo.Method
		exists[key] = true
	}

	for _, tt := range tests {
		key := tt.route + " " + tt.method
		if !exists[key] {
			t.Errorf("%s: not exists", tt.route)
		}
	}

}
