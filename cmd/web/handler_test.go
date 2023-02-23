package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
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

	pathToTemplates = "./../../templates/"

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

			req = addContexAndSessionToRequest(req, app)
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

func getCtx(req *http.Request) context.Context {
	ctx := context.WithValue(req.Context(), contextUserKey, "unknown")
	return ctx
}

func addContexAndSessionToRequest(req *http.Request, app application) *http.Request {
	req = req.WithContext(getCtx(req))

	ctx, _ := app.Session.Load(req.Context(), req.Header.Get("X-Session"))

	return req.WithContext(ctx)
}
