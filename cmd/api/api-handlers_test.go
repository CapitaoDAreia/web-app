package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
	"web-app/pkg/data"

	"github.com/go-chi/chi/v5"
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

func TestAppRefresh(t *testing.T) {
	var tests = []struct {
		name               string
		token              string
		expectedStatusCode int
		resetRefreshTime   bool
	}{
		{
			name:               "Valid",
			token:              "",
			expectedStatusCode: http.StatusOK,
			resetRefreshTime:   true,
		},
		{
			name:               "Valid but not yet ready to expire",
			token:              "",
			expectedStatusCode: http.StatusTooEarly,
			resetRefreshTime:   false,
		},
		{
			name:               "Expired token",
			token:              expiredToken,
			expectedStatusCode: http.StatusBadRequest,
			resetRefreshTime:   false,
		},
	}

	testUser := data.User{
		ID:        1,
		FirstName: "Admin",
		LastName:  "User",
		Email:     "admin@example.com",
	}

	oldRefreshTime := refreshTokenExpiry

	for _, test := range tests {
		var tkn string
		if test.token == "" {
			if test.resetRefreshTime {
				refreshTokenExpiry = time.Second * 1
			}
			tokens, _ := app.generateTokenPair(&testUser)
			tkn = tokens.RefreshToken
		} else {
			tkn = test.token
		}

		postedData := url.Values{
			"refresh_token": {tkn},
		}

		req, _ := http.NewRequest("POST", "/refresh-token", strings.NewReader(postedData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(app.refresh)
		handler.ServeHTTP(rr, req)

		if rr.Code != test.expectedStatusCode {
			t.Errorf("%s: expected status of %d but got %d", test.name, test.expectedStatusCode, rr.Code)
		}

		refreshTokenExpiry = oldRefreshTime
	}

}

func TestAppAllUserHandlers(t *testing.T) {
	var tests = []struct {
		name           string
		method         string
		json           string
		paramID        string
		handler        http.HandlerFunc
		expectedStatus int
	}{
		{"allUsers", "GET", "", "", app.allUsers, http.StatusOK},
		{"deleteUser", "DELETE", "", "1", app.deleteUser, http.StatusNoContent},
		{"deleteUser bad URL param", "DELETE", "", "Y", app.deleteUser, http.StatusBadRequest},
		{"getUser", "GET", "", "1", app.getUser, http.StatusOK},
		{"getUser invalid", "GET", "", "30", app.getUser, http.StatusBadRequest},
		{"getUser bad URL param", "GET", "", "Y", app.getUser, http.StatusBadRequest},

		{
			"updateUser valid",
			"PATCH",
			`{"id":1, "first_name": "Administrator", "last_name":"User", "email":"admin@example.com"}`,
			"",
			app.updateUser,
			http.StatusNoContent,
		},
		{
			"updateUser invalid",
			"PATCH",
			`{"id":1111, "first_name": "Administrator", "last_name":"User", "email":"admin@example.com"}`,
			"",
			app.updateUser,
			http.StatusBadRequest,
		},
		{
			"updateUser invalid json",
			"PATCH",
			`{"id":1, first_name: "Administrator", "last_name":"User", "email":"admin@example.com"}`,
			"",
			app.updateUser,
			http.StatusBadRequest,
		},
		{
			"insertUser valid",
			"PUT",
			`{"first_name": "Igor", "last_name":"Siilva", "email":"igor@example.com"}`,
			"",
			app.insertUser,
			http.StatusNoContent,
		},
		{
			"insertUser invalid",
			"PUT",
			`{"invalid":"invalid","first_name": "Igor", "last_name":"Siilva", "email":"igor@example.com"}`,
			"",
			app.insertUser,
			http.StatusBadRequest,
		},
		{
			"insertUser invalid json",
			"PUT",
			`{first_name: "Igor", "last_name":"Siilva", "email":"igor@example.com"}`,
			"",
			app.insertUser,
			http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		var req *http.Request

		if test.json == "" {
			req, _ = http.NewRequest(test.method, "/", nil)
		} else {
			req, _ = http.NewRequest(test.method, "/", strings.NewReader(test.json))
		}

		if test.paramID != "" {
			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("userID", test.paramID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(test.handler)

		handler.ServeHTTP(rr, req)

		if rr.Code != test.expectedStatus {
			t.Errorf("%s: wrong status returned; expected %d but got %d", test.name, test.expectedStatus, rr.Code)
		}
	}
}
