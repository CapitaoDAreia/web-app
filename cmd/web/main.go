package main

import (
	"log"
	"net/http"
)

type application struct{}

func main() {
	// Set up an app config
	app := application{}

	//get app routes
	mux := app.routes()

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
