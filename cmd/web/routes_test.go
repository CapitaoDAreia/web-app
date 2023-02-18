package main

import (
	"net/http"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestAppRoutes(t *testing.T) {
	var tests = []struct {
		route  string
		method string
	}{
		{"/", "GET"},
		{"/static/*", "GET"},
	}

	var app application
	mux := app.routes()

	chiRoutes := mux.(chi.Routes)

	for _, test := range tests {
		t.Run(test.route, func(t *testing.T) {
			if !routeExists(test.route, test.method, chiRoutes) {
				t.Errorf("route %s is note registered", test.route)
			}
		})
	}
}

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
