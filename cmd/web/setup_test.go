package main

import (
	"os"
	"testing"
)

var app application

func TestMain(m *testing.M) {
	//Execute before tests run
	os.Exit(m.Run())
}
