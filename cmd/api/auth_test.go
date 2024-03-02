package main

import (
	"24-Testing-Simple-Web-App/pkg/data"
	"fmt"
	"net/http/httptest"
	"testing"
)

func Test_application_getTokenFromHeaderAndVerify(t *testing.T) {

	testUser := data.User{
		ID:        1,
		FirstName: "Admin",
		LastName:  "User",
		Email:     "admin@example.com",
	}

	tokens, _ := app.generateTokenPairs(&testUser)

	var tests = []struct {
		name          string
		token         string
		errorExpected bool
		setHeader     bool
	}{
		{name: "valid", token: fmt.Sprintf("Bearer %s", tokens.Token), errorExpected: false, setHeader: true},
	}

	for _, tt := range tests {

		req := httptest.NewRequest("GET", "/", nil)
		if tt.setHeader {
			req.Header.Set("Authorization", tt.token)
		}

		rr := httptest.NewRecorder()

		_, _, err := app.getTokenFromHeaderAndVerify(rr, req)
		if err != nil {
			t.Errorf("%s: did not expect an eroro but got one %s", tt.name, err.Error())
		}
	}
}
