package main

import (
	"os"
	"testing"
	"web-app/pkg/repository/dbrepo"
)

var app application
var expiredToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiYXVkIjoiZXhhbXBsZS5jb20iLCJleHAiOjE2Nzc5NTkxOTcsImlzcyI6ImV4YW1wbGUuY29tIiwibmFtZSI6IkpvaG4gRG9lIiwic3ViIjoiMSJ9.NQquvm2FHC4uuH52iSgJT8qYZRORKqhFJ1Sln5QKKYE"

func TestMain(m *testing.M) {
	app.DB = &dbrepo.TestDBREpo{}
	app.Domain = "example.com"
	app.JWTSecret = "signin secret"
	os.Exit(m.Run())
}
