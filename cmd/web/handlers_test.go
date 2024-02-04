package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func Test_application_handlers(t *testing.T) {
	type args struct {
		url                string
		expectedStatusCode int
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "home", args: args{url: "/", expectedStatusCode: http.StatusOK}},
		{name: "404", args: args{url: "/finsh", expectedStatusCode: http.StatusNotFound}},
	}

	routes := app.routes()

	// create test web server
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	pathToTemplates = "./../../templates/"

	for _, tt := range tests {
		resp, err := ts.Client().Get(ts.URL + tt.args.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}
		if resp.StatusCode != tt.args.expectedStatusCode {
			t.Errorf("for %s: expected %d, but got %d", tt.name, tt.args.expectedStatusCode, resp.StatusCode)
		}
	}
}

func TestAppHome(t *testing.T) {
	var tests = []struct {
		name         string
		putInSession string
		expectedHTML string
		shouldFail   bool // flag to fail and pass
	}{
		{name: "firstVisit", putInSession: "", expectedHTML: "<small>From session:"},
		{name: "secondvisit", putInSession: "hello world!", expectedHTML: "<small>From session: hello world!"},
	}

	for _, tt := range tests {
		req := httptest.NewRequest("GET", "/", nil)

		req = addContextAndSessionToRequest(req, app)
		_ = app.Session.Destroy(req.Context())

		if tt.putInSession != "" {
			app.Session.Put(req.Context(), "test", tt.putInSession)
		}

		resp := httptest.NewRecorder()
		handler := http.HandlerFunc(app.Home)

		handler.ServeHTTP(resp, req)
		if resp.Code != http.StatusOK {
			t.Errorf("TestHome returned %d, but expected %d", resp.Code, http.StatusOK)
		}

		body, _ := io.ReadAll(resp.Body)

		if !strings.Contains(string(body), tt.expectedHTML) {
			t.Errorf("%s: did not find %s in response body ", tt.name, tt.expectedHTML)
		}
	}
}

func TestRenderBadTemplate(t *testing.T) {
	// set location to bad template
	pathToTemplates = "./testdata/"
	req := httptest.NewRequest("GET", "/", nil)
	req = addContextAndSessionToRequest(req, app)
	resp := httptest.NewRecorder()
	err := app.render(resp, req, "bad.page.gohtml", &TemplateData{})

	if err == nil {
		t.Error("Expected an error but not received, test failed")
	}

	pathToTemplates = "./../../templates/"
}

func getContext(req *http.Request) context.Context {
	ctx := context.WithValue(req.Context(), contextUserKey, "unknown")
	return ctx
}
func addContextAndSessionToRequest(req *http.Request, app application) *http.Request {
	req = req.WithContext(getContext(req))
	fmt.Println("Header: ", req.Header.Get("X-Session"))
	ctx, _ := app.Session.Load(req.Context(), req.Header.Get("X-Session"))

	return req.WithContext(ctx)
}

func Test_application_Login(t *testing.T) {
	type fields struct {
		postedData         url.Values
		expectedStatusCode int
		expectedLoc        string
	}

	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Valid login",
			fields: fields{
				postedData: url.Values{
					"email":    {"admin@example.com"},
					"password": {"secret"},
				},
				expectedStatusCode: http.StatusSeeOther,
				expectedLoc:        "/user/profile",
			},
		},
		{
			name: "Missing form data",
			fields: fields{
				postedData: url.Values{
					"email":    {""},
					"password": {""},
				},
				expectedStatusCode: http.StatusSeeOther,
				expectedLoc:        "/",
			},
		},
		{
			name: "User not found",
			fields: fields{
				postedData: url.Values{
					"email":    {"you@mail.com"},
					"password": {"password"},
				},
				expectedStatusCode: http.StatusSeeOther,
				expectedLoc:        "/",
			},
		},
		{
			name: "Bad credentials",
			fields: fields{
				postedData: url.Values{
					"email":    {"admin@example.com"},
					"password": {"password"},
				},
				expectedStatusCode: http.StatusSeeOther,
				expectedLoc:        "/",
			},
		},
	}
	for _, tt := range tests {
		req := httptest.NewRequest("POST", "/login", strings.NewReader(tt.fields.postedData.Encode()))
		req = addContextAndSessionToRequest(req, app)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(app.Login)
		handler.ServeHTTP(rr, req)

		if rr.Code != tt.fields.expectedStatusCode {
			t.Errorf("%s returned %d, but expected %d", tt.name, rr.Code, tt.fields.expectedStatusCode)
		}

		actualLoc, err := rr.Result().Location()
		if err == nil {
			if actualLoc.String() != tt.fields.expectedLoc {
				t.Errorf("%s: failed expected location is %s, but got %s", tt.name, actualLoc, tt.fields.expectedLoc)
			}
		} else {
			t.Errorf("%s: No location header set", tt.name)
		}

	}
}
