package main

import (
	"log"
	"os"
	"testing"
	"web-app/pkg/db"
)

var app application

func TestMain(m *testing.M) {
	pathToTemplates = "./../../templates/"

	app.Session = getSession()
	app.DataSourceName = `
	host=localhost 
	port=5432 
	user=postgres 
	password=postgres 
	dbname=users 
	timezone=UTC 
	connect_timeout=5`

	connection, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	app.DB = &db.PostgresConn{DB: connection}

	//Execute before tests run
	os.Exit(m.Run())
}
