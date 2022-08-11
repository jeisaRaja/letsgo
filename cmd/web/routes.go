package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.panicRecovery, app.logRequest, secureHeaders)
	dynamicMiddleware := alice.New(app.session.Enable)
	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/snippet/create", dynamicMiddleware.ThenFunc(app.create_snippet_form))
	mux.Post("/snippet/create", dynamicMiddleware.ThenFunc(app.create_snippet))
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.showSnippet))

	// User Registration and Login
	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signUpForm))
	fileserver := http.FileServer(http.Dir("./ui/static"))
	mux.Get("/static/", http.StripPrefix("/static", fileserver))
	return standardMiddleware.Then(mux)
}
