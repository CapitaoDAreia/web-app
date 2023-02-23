package main

import (
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
)

type application struct {
	Session *scs.SessionManager
}

func NewApplication() *application {
	return &application{}
}

func main() {
	// Set up an app config
	app := NewApplication()

	//get app routes
	mux := app.routes()

	//get a session manager
	app.Session = getSession()

	//print out a message
	log.Println("Server is listening on port 8080")

	// start the server
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

//manhã
//- explorando o painel para mapear as batidas na store-api e usar um client para efetuar esse tipo de operação
//- acompanhamos o Vini subindo a task de validação do controller de store-configurations da store-api --------- e aplicando um fix no way-to-go
