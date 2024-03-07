package main

import (
	"context"
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

func Test_application_UserHandlers(t *testing.T) {
	tests := []struct {
		name               string
		method             string
		json               string
		paramID            string
		handler            gin.HandlerFunc
		expectedStatusCode int
	}{
		{name: "AllUsers", method: "GET", json: "", paramID: "", handler: app.allUsers, expectedStatusCode: http.StatusOK},
		{name: "DeleteUser", method: "DELETE", json: "", paramID: "1", handler: app.deleteUser, expectedStatusCode: http.StatusNoContent},
		{name: "getUser valid", method: "GET", json: "", paramID: "1", handler: app.getUser, expectedStatusCode: http.StatusOK},
		{name: "getUser invalid", method: "GET", json: "", paramID: "100", handler: app.getUser, expectedStatusCode: http.StatusBadRequest},
	}

	for _, tt := range tests {
		var req *http.Request
		r := SetupServer()
		if tt.json == "" {
			switch tt.name {
			case "AllUsers":
				r.GET("/", tests[0].handler)
				return
			case "DeleteUser":
				r.DELETE("/", tests[1].handler)
				return
			case "getUser valid":
				r.GET("/", tests[2].handler)
				return
			case "etUser invalid":
				r.GET("/", tests[3].handler)
			}

			req = httptest.NewRequest(tt.method, "/", nil)
		} else {
			req = httptest.NewRequest(tt.method, "/", strings.NewReader(tt.json))
		}
		if tt.paramID != "" {
			ctx := gin.Context{}
			ctx.Params = gin.Params{
				{
					Key:   "userID",
					Value: tt.paramID,
				},
			}
			req = req.WithContext(context.WithValue(req.Context(), ctx.Param("userID"), ctx.Params))
		}
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		if rr.Code != tt.expectedStatusCode {
			t.Errorf("%s: failed! got %d, but expected %d", tt.name, rr.Code, tt.expectedStatusCode)
		}

	}
}

func SetupServer() *gin.Engine {
	r := gin.Default()
	gin.SetMode(gin.TestMode)
	return r
}
