package main

import (
	"os"
	"testing"
)

var app application

func TestMain(m *testing.M) {
	pathToTemplates = "./../../templates/"

	app.Session = getSession()

	//Execute before tests run
	os.Exit(m.Run())
}
