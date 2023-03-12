package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"web-app/pkg/data"
)

func TestAppEnableCORS(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})

	var tests = []struct {
		name           string
		method         string
		expectedHeader bool
	}{
		{
			name:           "preflight",
			method:         "OPTIONS",
			expectedHeader: true,
		},
		{
			name:           "get",
			method:         "GET",
			expectedHeader: false,
		},
	}

	for _, test := range tests {
		handlerToTest := app.enableCORS(nextHandler)

		req := httptest.NewRequest(test.method, "http://testing", nil)
		rr := httptest.NewRecorder()

		handlerToTest.ServeHTTP(rr, req)

		if test.expectedHeader && rr.Header().Get("Access-Control-Allow-Credentials") == "" {
			t.Errorf("%s: expected header, but did not find it", test.name)
		}

		if !test.expectedHeader && rr.Header().Get("Access-Control-Allow-Credentials") != "" {
			t.Errorf("%s: expected no header, but got one", test.name)
		}
	}
}

func TestAppAuthRequired(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	testUser := data.User{
		ID:        1,
		FirstName: "Admin",
		LastName:  "User",
		Email:     "admin@example.com",
	}

	tokens, _ := app.generateTokenPair(&testUser)

	var tests = []struct {
		name       string
		token      string
		expectAuth bool
		setHeader  bool
	}{
		{
			name:       "Valid token",
			token:      fmt.Sprintf("Bearer %s", tokens.Token),
			expectAuth: true,
			setHeader:  true,
		},
		{
			name:       "No token",
			token:      "",
			expectAuth: false,
			setHeader:  false,
		},
		{
			name:       "Invalid token",
			token:      fmt.Sprintf("Bearer %s", expiredToken),
			expectAuth: false,
			setHeader:  true,
		},
	}

	for _, test := range tests {
		req, _ := http.NewRequest("GET", "/", nil)
		if test.setHeader {
			req.Header.Set("Authorization", test.token)
		}

		rr := httptest.NewRecorder()

		handlerToTest := app.authRequired(nextHandler)

		handlerToTest.ServeHTTP(rr, req)

		if test.expectAuth && rr.Code == http.StatusUnauthorized {
			t.Errorf("%s: got code for 401, and should not have", test.name)
		}

		if !test.expectAuth && rr.Code != http.StatusUnauthorized {
			t.Errorf("%s: did not get code 401, and should have", test.name)
		}
	}
}
