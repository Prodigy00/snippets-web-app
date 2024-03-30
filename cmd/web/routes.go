package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux() // equivalent to express.Router()

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home) // eq to router.get("/", home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	return mux
}