package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAppAuthenticate(t *testing.T) {
	var tests = []struct {
		name               string
		requestBody        string
		expectedStatusCode int
	}{
		{
			name: "Valid user",
			requestBody: `{
				"email":"admin@example.com",
				"password":"secret"
			}`,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Invalid payload - Not json",
			requestBody:        `"emailadmin@example.com","passwordsecret"`,
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:               "Invalid payload - Empty json",
			requestBody:        "{}",
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name: "No password",
			requestBody: `{
				"email":"admin@example.com",
				"password":""
			}`,
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name: "No email",
			requestBody: `{
				"email":"",
				"password":"secret"
			}`,
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name: "Invalid user",
			requestBody: `{
				"email":"admin@example1111.com",
				"password":"secret"
			}`,
			expectedStatusCode: http.StatusUnauthorized,
		},
	}

	for _, test := range tests {
		var reader io.Reader
		reader = strings.NewReader(test.requestBody)
		req, _ := http.NewRequest("POST", "/auth", reader)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(app.authenticate)

		handler.ServeHTTP(rr, req)

		if test.expectedStatusCode != rr.Code {
			t.Errorf("%s: returned wrong status code; expected %d but got %d", test.name, test.expectedStatusCode, rr.Code)
		}
	}
}
