package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	chain "github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	//initialize the router
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})
	//Leave the static files route unchanged No need for a session	
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	dynamic := chain.New(app.sessionManager.LoadAndSave)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet, "/snippet/create", dynamic.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", dynamic.ThenFunc(app.snippetCreatePost))

	// "/static" prefix before the request reaches the file server.

	// return app.recoverPanic(app.logRequest(secureHeaders(router)))
	standard := chain.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(router)
}
