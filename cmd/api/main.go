package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"web-app/pkg/repository"
	"web-app/pkg/repository/dbrepo"
)

const PORT = 8090

type application struct {
	DSN       string
	DB        repository.DatabaseRepo
	Domain    string
	JWTSecret string
}

func main() {
	var app application
	flag.StringVar(&app.Domain, "domain", "example.com", "Domain for appliation, e.g. company.com")
	flag.StringVar(
		&app.DSN,
		"DSN",
		`host=localhost 
		port=5432 
		user=postgres 
		password=postgres 
		dbname=users 
		timezone=UTC 
		connect_timeout=5`,
		"Postgres connection",
	)
	flag.StringVar(&app.JWTSecret, "jwt-secret", "secret", "signin secret")
	flag.Parse()

	connection, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	app.DB = &dbrepo.PostgresDBRepo{DB: connection}

	log.Printf("Server is listenin on PORT %d\n", PORT)

	err = http.ListenAndServe(fmt.Sprintf(":%d", PORT), app.routes())
}
