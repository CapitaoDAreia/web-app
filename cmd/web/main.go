package main

import (
	"encoding/gob"
	"flag"
	"log"
	"net/http"
	"web-app/pkg/data"
	"web-app/pkg/db"

	"github.com/alexedwards/scs/v2"
)

type application struct {
	DataSourceName string
	DB             *db.PostgresConn
	Session        *scs.SessionManager
}

func NewApplication() *application {
	return &application{}
}

func main() {
	gob.Register(data.User{})
	// Set up an app config
	app := NewApplication()

	flag.StringVar(
		&app.DataSourceName,
		"DataSourceName",
		`host=localhost 
		port=5432 
		user=postgres 
		password=postgres 
		dbname=users 
		timezone=UTC 
		connect_timeout=5`,
		"Postgres connection",
	)
	flag.Parse()

	connection, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	app.DB = &db.PostgresConn{DB: connection}

	//get a session manager
	app.Session = getSession()

	//get app routes
	mux := app.routes()

	//print out a message
	log.Println("Server is listening on port 8080")

	// start the server
	if err = http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

//manhã
//- explorando o painel para mapear as batidas na store-api e usar um client para efetuar esse tipo de operação
//- acompanhamos o Vini subindo a task de validação do controller de store-configurations da store-api --------- e aplicando um fix no way-to-go
