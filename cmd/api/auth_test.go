package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"web-app/pkg/data"
)

func TestAppGetTokenFromHeaderAndVerify(t *testing.T) {
	testUser := data.User{
		ID:        1,
		FirstName: "Admin",
		LastName:  "User",
		Email:     "admin@example.com",
	}

	tokens, _ := app.generateTokenPair(&testUser)

	var tests = []struct {
		name          string
		token         string
		errorExpected bool
		setHeader     bool
		issuer        string
	}{
		{
			name:          "Valid",
			token:         fmt.Sprintf("Bearer %s", tokens.Token),
			errorExpected: false,
			setHeader:     true,
			issuer:        app.Domain,
		},
		{
			name:          "Valid but expired",
			token:         fmt.Sprintf("Bearer %s", expiredToken),
			errorExpected: true,
			setHeader:     true,
			issuer:        app.Domain,
		},
		{
			name:          "No header",
			token:         "",
			errorExpected: true,
			setHeader:     false,
			issuer:        app.Domain,
		},
		{
			name:          "Invalid token",
			token:         fmt.Sprintf("Bearer %sinvalidToken", tokens.Token),
			errorExpected: true,
			setHeader:     true,
			issuer:        app.Domain,
		},
		{
			name:          "No bearer",
			token:         fmt.Sprintf("Bear %s", tokens.Token),
			errorExpected: true,
			setHeader:     true,
			issuer:        app.Domain,
		},
		{
			name:          "Three header parts",
			token:         fmt.Sprintf("Bearer %s excedent", tokens.Token),
			errorExpected: true,
			setHeader:     true,
			issuer:        app.Domain,
		},
		// because of the appDomain in the test loop, this must be the last case that will run
		{
			name:          "Wrong issuer",
			token:         fmt.Sprintf("Bearer %s", tokens.Token),
			errorExpected: true,
			setHeader:     true,
			issuer:        "wrongDomain.com",
		},
	}

	for _, test := range tests {
		if test.issuer != app.Domain {
			app.Domain = test.issuer
			tokens, _ = app.generateTokenPair(&testUser)
		}
		req, _ := http.NewRequest("GET", "/", nil)
		if test.setHeader {
			req.Header.Set("Authorization", test.token)
		}

		rr := httptest.NewRecorder()

		_, _, err := app.getTokenFromHeaderAndVerify(rr, req)
		if err != nil && !test.errorExpected {
			t.Errorf("%s: did not expect error, but got one - %s", test.name, err.Error())
		}

		if err == nil && test.errorExpected {
			t.Errorf("%s: expected error, but did not get one", test.name)
		}

		app.Domain = "example.com"
	}
}
