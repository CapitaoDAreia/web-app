package main

import (
	"os"
	"testing"
	"web-app/pkg/repository/dbrepo"
)

var app application

func TestMain(m *testing.M) {
	pathToTemplates = "./../../templates/"

	app.Session = getSession()

	app.DB = &dbrepo.TestDBREpo{}

	//Execute before tests run
	os.Exit(m.Run())
}
