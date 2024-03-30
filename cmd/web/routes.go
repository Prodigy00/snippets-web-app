package main

import (
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux() // equivalent to express.Router()

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home) // eq to router.get("/", home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// Pass the servemux as the 'next' parameter to the secureHeaders middleware.
	//Because secureHeaders is just a function, and the function returns a
	// http.Handler we don't need to do anything else.
	return standard.Then(mux)
}
