package main

import (
	"24-Testing-Simple-Web-App/pkg/data"
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"image"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
	"testing"
)

func Test_application_handlers(t *testing.T) {
	type args struct {
		url                     string
		expectedStatusCode      int
		expectedURL             string
		expectedFirstStatusCode int
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "home", args: args{url: "/",
			expectedStatusCode:      http.StatusOK,
			expectedURL:             "/",
			expectedFirstStatusCode: http.StatusOK}},
		{name: "404",
			args: args{
				url:                     "/finsh",
				expectedStatusCode:      http.StatusNotFound,
				expectedURL:             "/finsh",
				expectedFirstStatusCode: http.StatusNotFound},
		},
		{
			name: "profile",
			args: args{
				url:                     "/user/profile",
				expectedStatusCode:      http.StatusOK,
				expectedURL:             "/",
				expectedFirstStatusCode: http.StatusTemporaryRedirect,
			},
		},
	}

	routes := app.routes()

	// create test web server
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{
		Transport: tr,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	for _, tt := range tests {
		resp, err := ts.Client().Get(ts.URL + tt.args.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}
		if resp.StatusCode != tt.args.expectedStatusCode {
			t.Errorf("for %s: expected %d, but got %d", tt.name, tt.args.expectedStatusCode, resp.StatusCode)
		}
		if resp.Request.URL.Path != tt.args.expectedURL {
			t.Errorf("%s: expected final url of %s, but got %s", tt.name, tt.args.expectedURL, resp.Request.URL.Path)
		}

		resp2, _ := client.Get(ts.URL + tt.args.url)
		if resp2.StatusCode != tt.args.expectedFirstStatusCode {
			t.Errorf("%s: expected first status code %d, but got %d", tt.name, tt.args.expectedFirstStatusCode, resp2.StatusCode)
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

func Test_application_UploadFiles(t *testing.T) {
	pr, pw := io.Pipe()

	writer := multipart.NewWriter(pw)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go simulatePNGUpload("./testdata/img.png", *writer, t, wg)

	request := httptest.NewRequest("POST", "/", pr)
	request.Header.Add("Content-Type", writer.FormDataContentType())

	uploadedFiles, err := app.UploadFiles(request, "./testdata/uploads/")
	if err != nil {
		t.Error(err)
	}

	if _, err := os.Stat(fmt.Sprintf("./testdata/uploads/%s", uploadedFiles[0].OriginalFileName)); os.IsNotExist(err) {
		t.Errorf("expect file to exists %s", err.Error())
	}

	_ = os.Remove(fmt.Sprintf("./testdata/uploads/%s", uploadedFiles[0].OriginalFileName))

}

func simulatePNGUpload(fileToUpload string, writer multipart.Writer, t *testing.T, wg *sync.WaitGroup) {
	defer func(writer *multipart.Writer) {
		err := writer.Close()
		if err != nil {
			return
		}
	}(&writer)

	defer wg.Done()
	part, err := writer.CreateFormFile("file", path.Base(fileToUpload))
	if err != nil {
		t.Error(err)
	}

	f, err := os.Open(fileToUpload)
	if err != nil {
		t.Error(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			return
		}
	}(f)

	img, _, err := image.Decode(f)
	if err != nil {
		t.Error(err)
	}

	err = png.Encode(part, img)
	if err != nil {
		t.Error(err)
	}
}

func Test_application_UploadProfilePicture(t *testing.T) {
	uploadPath = "./testdata/uploads"
	filePath := "./testdata/img.png"

	// field name for form
	fieldName := "file"

	// bytes.Buffer as request body
	body := new(bytes.Buffer)

	mw := multipart.NewWriter(body)
	file, err := os.Open(filePath)
	if err != nil {
		t.Errorf("file can not open %s", err.Error())
	}

	w, err := mw.CreateFormFile(fieldName, filePath)
	if err != nil {
		t.Errorf("file can not open %s", err.Error())
	}

	if _, err := io.Copy(w, file); err != nil {
		t.Errorf(err.Error())
	}

	err = mw.Close()
	if err != nil {
		return
	}

	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	req = addContextAndSessionToRequest(req, app)
	app.Session.Put(req.Context(), "user", data.User{ID: 1})
	req.Header.Add("Content-Type", mw.FormDataContentType())

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.UploadProfilePicture)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("wrong status code %d want: %d", rr.Code, http.StatusSeeOther)
	}

	_ = os.Remove("./testdata/uploads/img.png")
}
