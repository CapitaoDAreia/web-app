package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIPFromContext(t *testing.T) {
	//create an app of type application
	var app application

	//get a context
	ctx := context.Background()

	// put something in the context
	ctx = context.WithValue(ctx, contextUserKey, "anyValue")
	// call function
	ip := app.ipFromContext(ctx)

	if !strings.EqualFold("anyValue", ip) {
		t.Errorf("Incorrect context value. ctx: %v, ip: %v", ctx, ip)
	}
}

func TestAppAddIPToContext(t *testing.T) {
	tests := []struct {
		testName    string
		headerName  string
		headerValue string
		addr        string
		emptyAddr   bool
	}{
		{
			testName:    "default request",
			headerName:  "",
			headerValue: "",
			addr:        "",
			emptyAddr:   false,
		},
		{
			testName:    "default request with empty addr",
			headerName:  "",
			headerValue: "",
			addr:        "",
			emptyAddr:   true,
		},
		{
			testName:    "exists an header",
			headerName:  "X-Forwarded-For",
			headerValue: "::1",
			addr:        "",
			emptyAddr:   false,
		},
		{
			testName:    "invalid address",
			headerName:  "",
			headerValue: "",
			addr:        "hello:world",
			emptyAddr:   false,
		},
		// {
		// 	testName:    "",
		// 	headerName:  "",
		// 	headerValue: "",
		// 	addr:        "",
		// 	emptyAddr:   false,
		// },
	}

	// dummy handler
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//make sure that the value exists in the context
		val := r.Context().Value(contextUserKey)
		if val == nil {
			t.Error(contextUserKey, "not present")
		}

		// make sure we got a string back
		ip, ok := val.(string)
		if !ok {
			t.Error("not string")
		}
		t.Log(ip)
	})

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// create handler to test
			handlerToTest := app.addIPToContext(nextHandler)

			req := httptest.NewRequest("GET", "http://testing", nil)

			if test.emptyAddr {
				req.RemoteAddr = ""
			}

			if len(test.headerName) > 0 {
				req.Header.Add(test.headerName, test.headerValue)
			}

			if len(test.addr) > 0 {
				req.RemoteAddr = test.addr
			}

			handlerToTest.ServeHTTP(httptest.NewRecorder(), req)
		})
	}
}
