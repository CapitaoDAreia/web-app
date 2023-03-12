package main

import (
	"net/http"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestAppRoutes(t *testing.T) {
	//create test table with cases
	var tests = []struct {
		testName string
		route    string
		method   string
	}{
		{"success on auth", "/auth", "POST"},
		{"success on refresh-token", "/refresh-token", "POST"},
		{"", "/refresh-token", "POST"},
		{"", "/users/", "GET"},
		{"", "/users/{userID}", "GET"},
		{"", "/users/{userID}", "DELETE"},
		{"", "/users/", "PATCH"},
		{"", "/users/", "PUT"},
	}

	// creates an http.Handler
	mux := app.routes()

	//casting the handler to 'chi.Routes' type
	chiRoutes := mux.(chi.Routes)

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			if !routeExists(test.route, test.method, chiRoutes) {
				t.Errorf("route %s is note registered", test.route)
			}
		})
	}
}

// validatte if route exists based on chi.walk function
func routeExists(testRoute, testMethod string, chiRoutes chi.Routes) bool {
	found := false

	_ = chi.Walk(chiRoutes, func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if strings.EqualFold(method, testMethod) && strings.EqualFold(route, testRoute) {
			found = true
		}
		return nil
	})

	return found
}
