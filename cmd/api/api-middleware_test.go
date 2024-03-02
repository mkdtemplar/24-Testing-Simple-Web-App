package main

import (
	"24-Testing-Simple-Web-App/pkg/data"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_application_authRequired(t *testing.T) {
	nextHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

	})

	testUser := data.User{
		ID:        1,
		FirstName: "Admin",
		LastName:  "User",
		Email:     "admin@example.com",
	}

	tokens, _ := app.generateTokenPairs(&testUser)

	tests := []struct {
		name             string
		token            string
		expectAuthorized bool
		setHeader        bool
	}{
		{name: "valid token", token: fmt.Sprintf("Bearer %s", tokens.Token), expectAuthorized: true, setHeader: true},
	}
	for _, tt := range tests {
		req := httptest.NewRequest("GET", "/", nil)
		if tt.setHeader {
			req.Header.Set("Authorization", tt.token)
		}

		rr := httptest.NewRecorder()
		handlerToTest := app.authRequired(nextHandler)
		handlerToTest.ServeHTTP(rr, req)

		if tt.expectAuthorized && rr.Code == http.StatusUnauthorized {
			t.Errorf("%s: got code %d but should have %d", tt.name, http.StatusUnauthorized, http.StatusOK)
		}

		if !tt.expectAuthorized && rr.Code != http.StatusUnauthorized {
			t.Errorf("%s: did not get code %d but should have %d", tt.name, http.StatusUnauthorized, http.StatusUnauthorized)
		}
	}
}
