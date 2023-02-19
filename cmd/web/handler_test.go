package main

import (
	"io"
	"net/http/httptest"
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

	var app application
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
