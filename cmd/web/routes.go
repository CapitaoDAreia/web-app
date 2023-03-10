package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(app.addIPToContext)
	mux.Use(app.Session.LoadAndSave)

	mux.Get("/", app.HomeTemplate)
	mux.Get("/login", app.LoginTemplate)

	mux.Post("/login", app.Login)

	mux.Get("/user/profile", app.ProfileTemplate)

	mux.Route("/user", func(mux chi.Router) {
		mux.Use(app.auth)
		mux.Get("/profile", app.ProfileTemplate)
	})

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
