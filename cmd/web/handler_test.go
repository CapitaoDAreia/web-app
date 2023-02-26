package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// create an test function
func TestRender(t *testing.T) {
	// template, _ := os.ReadFile("../../templates/home.page.gohtml")
	// expectedTemplate := bytes.NewReader(template)

	tests := []struct {
		testName             string
		URL                  string
		expectedStatusCode   int
		expectedResponseBody io.Reader
	}{
		{
			testName: "Home",
			URL:      "/", expectedStatusCode: 200,
			// expectedResponseBody: expectedTemplate,
		},
		{
			testName:           "Home",
			URL:                "/",
			expectedStatusCode: 200,
			// expectedResponseBody: expectedTemplate,
		},
	}

	routes := app.routes()

	server := httptest.NewTLSServer(routes)
	defer server.Close()

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// fmt.Println(server.URL)
			// fmt.Println(test.URL)
			serverResponse, err := server.Client().Get(server.URL + test.URL)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if serverResponse.StatusCode != test.expectedStatusCode {
				t.Errorf("Error, status unexpected! Expected %v, have %v.", test.expectedStatusCode, serverResponse.StatusCode)
			}

			if serverResponse.Body == nil {
				t.Errorf("body is nil")
			}
		})
	}
}

func TestAppHome(t *testing.T) {
	tests := []struct {
		name         string
		putInSession string
		expectedHTML string
	}{
		{
			name:         "Success on first visit",
			putInSession: "",
			expectedHTML: "<p>From session:",
		},
		{
			name:         "Success on second visit",
			putInSession: "hello, world!",
			expectedHTML: "<p>From session: hello, world!",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//Create a request
			req, _ := http.NewRequest("GET", "/", nil)

			req = addContextAndSessionToRequest(req, app)
			_ = app.Session.Destroy(req.Context())

			if test.putInSession != "" {
				app.Session.Put(req.Context(), "test", test.putInSession)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(app.HomeTemplate)
			handler.ServeHTTP(rr, req)

			//Check status code
			if rr.Code != http.StatusOK {
				t.Errorf("TestAppHome: expected 200 but got %v", rr.Code)
			}

			body, _ := io.ReadAll(rr.Body)
			if !strings.Contains(string(body), test.expectedHTML) {
				t.Errorf("%s: did not find %s in response body", test.name, test.expectedHTML)
			}
		})
	}
}

func TestAppRenderWithBadTemplate(t *testing.T) {
	//set template path (pathToTemplates) to a location with a bad template
	pathToTemplates = "./testdata/"

	req, _ := http.NewRequest("GET", "/", nil)
	req = addContextAndSessionToRequest(req, app)
	rr := httptest.NewRecorder()

	err := app.render(rr, req, "bad.page.gohtml", &TemplateData{})
	if err == nil {
		t.Error("Expected error from bad template, but did not get one.")
	}

	pathToTemplates = "./../../templates/"
}

func getCtx(req *http.Request) context.Context {
	ctx := context.WithValue(req.Context(), contextUserKey, "unknown")
	return ctx
}

func addContextAndSessionToRequest(req *http.Request, app application) *http.Request {
	req = req.WithContext(getCtx(req))

	ctx, _ := app.Session.Load(req.Context(), req.Header.Get("X-Session"))

	return req.WithContext(ctx)
}

func TestAppLogin(t *testing.T) {
	tests := []struct {
		name               string
		postedData         url.Values
		expectedStatusCode int
		expectedLocation   string
	}{
		{
			name: "Valid login",
			postedData: url.Values{
				"email":    {"admin@example.com"},
				"password": {"secret"},
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLocation:   "/user/profile",
		},
		{
			name: "Missing form data",
			postedData: url.Values{
				"email":    {""},
				"password": {""},
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLocation:   "/",
		},
		{
			name: "User not found",
			postedData: url.Values{
				"email":    {"invalid@example.com"},
				"password": {"invalid"},
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLocation:   "/",
		},
		{
			name: "Bad credentials",
			postedData: url.Values{
				"email":    {"admin@example.com"},
				"password": {"wrong"},
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLocation:   "/",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/login", strings.NewReader(test.postedData.Encode()))
			req = addContextAndSessionToRequest(req, app)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(app.Login)
			handler.ServeHTTP(rr, req)

			if rr.Code != test.expectedStatusCode {
				t.Errorf("Expected status code %v but got %v", test.expectedStatusCode, rr.Code)
			}

			actualLoc, err := rr.Result().Location()
			if err == nil {
				if actualLoc.String() != test.expectedLocation {
					t.Errorf("Expected location to be '%v' but got '%v'", test.expectedLocation, actualLoc.String())
				}
			} else {
				t.Errorf("No location header set.")
			}

		})
	}
}
