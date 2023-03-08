package main

import (
	"os"
	"testing"
	"web-app/pkg/repository/dbrepo"
)

var app application

func TestMain(m *testing.M) {
	app.DB = &dbrepo.TestDBREpo{}
	app.Domain = "example.com"
	app.JWTSecret = "signin secret"
	os.Exit(m.Run())
}
