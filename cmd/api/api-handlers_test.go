package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func Test_application_authenticate(t *testing.T) {

	type args struct {
		name               string
		requestBody        string
		expectedStatusCode int
	}
	tests := []struct {
		args args
	}{
		{
			args: args{
				name:               "valid",
				requestBody:        `{"email": "admin@example.com", "password": "secret"}`,
				expectedStatusCode: 200,
			},
		},
		{
			args: args{
				name:               "not JSON",
				requestBody:        "Bad JSON",
				expectedStatusCode: http.StatusBadRequest,
			},
		},
		{
			args: args{
				name:               "Empty JSON",
				requestBody:        `{}`,
				expectedStatusCode: http.StatusUnauthorized,
			},
		},
	}
	for _, tt := range tests {
		var reader io.Reader
		reader = strings.NewReader(tt.args.requestBody)
		r := SetupServer()
		r.POST("/auth", app.authenticate)

		req := httptest.NewRequest("POST", "/auth", reader)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		if tt.args.expectedStatusCode != resp.Code {
			t.Errorf("expected code %d, but got %d", tt.args.expectedStatusCode, resp.Code)
		}
	}
}

func SetupServer() *gin.Engine {
	r := gin.Default()
	gin.SetMode(gin.TestMode)
	return r
}
