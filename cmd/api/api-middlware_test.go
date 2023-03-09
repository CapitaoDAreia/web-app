package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
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
