package main

import (
	"reflect"
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

	var output []struct {
		route  string
		method string
	}
	mux := app.routes()

	for _, routeInfo := range mux.Routes() {
		for _, tt := range tests {
			if routeInfo.Path == tt.route && routeInfo.Method == tt.method {
				output = append(output, struct {
					route  string
					method string
				}{route: routeInfo.Path, method: routeInfo.Method})
			}
		}
	}

	if reflect.DeepEqual(tests, output) == false {
		t.Errorf("%s: failed", "routes()")
	}
}
